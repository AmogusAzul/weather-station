#ifndef SD_HANDLER_H
#define SD_HANDLER_H

#include <cstdint>
#include <SdFat.h>
#include "../common/types/base32-bytes-types/base32-bytes-types.h"
#include "../communication-types/byte-types/communication-bytes-types.h"
#include "../common/types/string-bytes-types/string-bytes-types.h"
#include "../safety/safety.h"
#include "../communication-types/struct-likes/struct-likes.h"
#include "../common/types/bytes-types.h"

#include <SPI.h>

// Ensure ENABLE_DEDICATED_SPI is only defined once
#ifndef ENABLE_DEDICATED_SPI
#define ENABLE_DEDICATED_SPI 0
#endif

// SD_FAT_TYPE = 0 for SdFat/File as defined in SdFatConfig.h,
// 1 for FAT16/FAT32, 2 for exFAT, 3 for FAT16/FAT32 and exFAT.
#define SD_FAT_TYPE 0

// Define the SD and file types based on SD_FAT_TYPE
#if SD_FAT_TYPE == 0
using SdType = SdFat;
using FileType = SdFile;
#elif SD_FAT_TYPE == 1
using SdType = SdFat32;
using FileType = File32;
#elif SD_FAT_TYPE == 2
using SdType = SdExFat;
using FileType = ExFile;
#elif SD_FAT_TYPE == 3
using SdType = SdFs;
using FileType = FsFile;
#else
#error "Invalid SD_FAT_TYPE"
#endif

// Define file names as constants
const char* const LOG_FILE = "log.txt";
const char* const TOKEN_FILE = "token.bin";
const char* const BACKUP_INFO_FILE = "backup-info.bin";
const char* const BACKUP_FORMAT_FILE = "backup-format.bin";
const char* const NETWORK_CONFIG_FILE = "network-config.bin";

template <uint8_t misoPin, uint8_t mosiPin, uint8_t sckPin>
class SDHandler {
private:
    // Pins and driver for SoftSPI
    const uint8_t csPin;
    SoftSpiDriver<misoPin, mosiPin, sckPin> softSpi;  // SoftSpiDriver instantiated with template parameters
    SdType sd;
    FileType file;
    bool sdInitialized;  // Stores the result of SD initialization

public:
    // Constructor with template-based pins
    SDHandler(const uint8_t csPin) : csPin(csPin), softSpi(), sdInitialized(false) {
        // Configure SD using SdSpiConfig
        SdSpiConfig config(csPin, SHARED_SPI, SD_SCK_MHZ(0), &softSpi);
        sdInitialized = sd.begin(config);
    }

    // Destructor to clean up
    ~SDHandler() = default;

    void logToFile(const std::string& logEntry) {
        // Open the log file in append mode (O_APPEND)
        if (!file.open(LOG_FILE, O_WRONLY | O_CREAT | O_APPEND)) {
            return;
        }


        // Write the log entry followed by a newline
        file.println(logEntry.c_str());

        file.close();  // Close the file
    }

        // Write the token in the file (swap the current token with a new one)
    void swapToken(const TokenBytes &newToken) {
        if (!sdInitialized) return;  // Check if SD card is initialized

        // Open the token file
        if (!file.open(TOKEN_FILE, O_RDWR | O_CREAT)) {
            return;
        }

        // Write the new token
        file.rewind();
        file.write(&safety::TOKEN_LENGTH, 1);  // Write the new token length
        file.write(newToken.getBytes(), safety::TOKEN_LENGTH);  // Write the new token

        file.close();  // Close the file
    }

    // Read the current token from the file and return a new TokenBytes object
    TokenBytes readToken() {

        // If SD card is not initialized or file cannot be opened, return default token
        if (!sdInitialized || !file.open(TOKEN_FILE, O_RDONLY)) {
            return TokenBytes(safety::DEFAULT_TOKEN);
        }

        // Check if the file length matches TOKEN_LENGTH + 1 (for the length byte)
        if (file.fileSize() != safety::TOKEN_LENGTH + 1) {
            file.close();
            return TokenBytes(safety::DEFAULT_TOKEN);
        }

        //skip TOKEN_LENGTH byte
        file.seekSet(1);

        // Read all the token bytes
        uint8_t buffer[safety::TOKEN_LENGTH];
        
        // Ensure that TOKEN_LENGTH bytes are read successfully
        if (file.read(buffer, safety::TOKEN_LENGTH) != safety::TOKEN_LENGTH) {
            file.close();
            return TokenBytes(safety::DEFAULT_TOKEN);  // Return default token in case of read error
        }

        file.close();  // Close the file

        // Return the new TokenBytes object with the data
        return TokenBytes(buffer);
    }

    void writeBackupInfo(const Packet packet) {
        // Open the backup info file in append mode (O_APPEND)
        if (!file.open(BACKUP_INFO_FILE, O_WRONLY | O_CREAT | O_APPEND)) {
            return;  // Exit if file couldn't be opened
        }

        // Write the backup data
        file.write(packet.getBytes(), packet.getLength());

        file.close();  // Close the file

         // Now open the FORMAT file in write mode (overwrite existing content)
        if (!file.open(BACKUP_FORMAT_FILE, O_WRONLY | O_CREAT | O_TRUNC)) {
            return;  // Exit if file couldn't be opened
        }

        Bytes format = packet.getFormat();

        // Write the hardcoded byte array to the file
        file.write(format.getBytes(), format.getLength());

        // Close the file
        file.close();
    }

    size_t getBackupLength() {
        file.open(BACKUP_INFO_FILE, O_RDONLY);
        size_t infoSize = file.fileSize();
        file.close();

        file.open(BACKUP_FORMAT_FILE, O_RDONLY);
        size_t lengthPerEntry;
        //Skipping version
        file.seekSet(1);    
        // Buffer for two bytes
        uint8_t buffer[2];

        // Read two bytes at a time
        for (size_t i = 2; i < file.fileSize(); i += 2) {
            if (file.read(buffer, 2) == 2) {
                lengthPerEntry += buffer[1];  // We only need the second byte (length)
            } else {
                // Handle read error (optional)
                break;
            }
        }

        return infoSize/lengthPerEntry;

    }

    // Read the network configuration from the file and return a NetworkData object
    NetworkData readNetworkConfig() {
        Serial.println("start function");

        NetworkData networkData;  // Default network data

        Serial.println("made base NetworkData object");

        if (!sdInitialized || !file.open(NETWORK_CONFIG_FILE, O_RDONLY)) {

            Serial.println("error while opening file");
            return networkData;  // Return default values if SD card isn't initialized or file can't be opened
        }

        Serial.println("opened file successfully");

        //reading all bytes
        size_t fileSize = file.fileSize();
        uint8_t buffer[fileSize];

        if (file.read(buffer, fileSize) != fileSize) {
            Serial.println("couldn't read file");
            file.close();
            return networkData;
        }
        file.close();

        Serial.println("read file successfully");

        // After reading the file
        Serial.println("File read successfully. Printing raw content:");
        for (size_t i = 0; i < fileSize; ++i) {
            Serial.print((char)buffer[i]);
        }
        Serial.println(); // New line after printing raw content

        size_t offset;
        // Debug for Station ID parsing
        Serial.println("Parsing Station ID...");
        UInt32Bytes stationId(buffer[offset]);  // Store station ID as UInt32Bytes
        Serial.print("Station ID: "); Serial.println(stationId.getNumber());
        offset += stationId.getLength();

        // Debug for SSID parsing
        Serial.println("Parsing SSID...");
        StringBytes ssid(buffer[offset] + 1, &buffer[offset]);
        Serial.print("SSID: "); Serial.println(ssid.getString().c_str());
        offset += ssid.getLength();

        // Debug for Password parsing
        Serial.println("Parsing Password...");
        StringBytes password(buffer[offset] + 1, &buffer[offset]);
        Serial.print("Password: "); Serial.println(password.getString().c_str());
        offset += password.getLength();

        // Debug for IP parsing
        Serial.println("Parsing IP...");
        StringBytes ip(static_cast<size_t>(buffer[offset] + 1), &buffer[offset]);
        Serial.print("IP: "); Serial.println(ip.getString().c_str());
        offset += ip.getLength();

        // Debug for Port parsing
        Serial.println("Parsing Port...");
        UInt32Bytes port(buffer);
        Serial.print("Port: "); Serial.println(port.getNumber());
        offset += port.getLength();

        return NetworkData(stationId, ip, port, ssid, password);  // Return populated NetworkData object
    };

    bool isSDInitialized() const { return sdInitialized; };

};

#endif // SD_HANDLER_H

