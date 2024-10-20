#ifndef SD_HANDLER_H
#define SD_HANDLER_H

#include <cstdint>

#include <SdFat.h>
#include "../../common/types/base32-bytes-types/base32-bytes-types.h"
#include "../communication-types/byte-types/communication-bytes-types.h"
#include "../../common/types/string-bytes-types/string-bytes-types.h"
#include "../safety/safety.h"
#include "../communication-types/struct-likes/struct-likes.h"

// Define file names as constants
const char* const LOG_FILE = "log.txt";
const char* const TOKEN_FILE = "tokens.bin";
const char* const BACKUP_INFO_FILE = "backup-info.bin";
const char* const BACKUP_FORMAT_FILE = "backup-format.bin";
const char* const NETWORK_CONFIG_FILE = "network-config.bin";

template <uint8_t misoPin, uint8_t mosiPin, uint8_t sckPin>
class SDHandler {
public:
    // Constructor with template-based pins
    SDHandler(const uint8_t csPin) : csPin(csPin), sdInitialized(false) {
        // Initialize SoftSPI driver at compile time using template parameters
        softSpi = new SoftSpiDriver<misoPin, mosiPin, sckPin>();

        // SD card initialization using the template-driven SoftSPI driver
        SdSpiConfig config(csPin, SHARED_SPI, SD_SCK_MHZ(0), softSpi);
        sdInitialized = sd.begin(config);
    }

    // Destructor to clean up
    ~SDHandler() {
        delete softSpi;
    }

    void SDHandler::logToFile(const std::string& logEntry) {
        // Open the log file in append mode (O_APPEND)
        File file = sd.open(LOG_FILE, O_WRONLY | O_CREAT | O_APPEND);
        if (!file) {
            return;  // Exit if file couldn't be opened
        }

        // Write the log entry followed by a newline
        file.println(logEntry.c_str());

        file.close();  // Close the file
    }

        // Write the token in the file (swap the current token with a new one)
    void SDHandler::swapToken(const TokenBytes &newToken) {
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
    TokenBytes SDHandler::readToken() {

        // If SD card is not initialized or file cannot be opened, return default token
        if (!sdInitialized || !file.open(TOKEN_FILE, O_RDONLY)) {
            return TokenBytes(&(safety::DEFAULT_TOKEN));
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

    void SDHandler::writeBackupInfo(const ByteArray& backupData) {
        // Open the backup info file in append mode (O_APPEND)
        File file = sd.open(BACKUP_INFO_FILE, O_WRONLY | O_CREAT | O_APPEND);
        if (!file) {
            return;  // Exit if file couldn't be opened
        }

        // Write the backup data
        file.write(backupData.getBytes(), backupData.getLength());

        file.close();  // Close the file
    }

    // Read the network configuration from the file and return a NetworkData object
    NetworkData SDHandler::readNetworkConfig() {
        NetworkData networkData;  // Default network data

        if (!sdInitialized || !file.open(NETWORK_CONFIG_FILE, O_RDONLY)) {
            return networkData;  // Return default values if SD card isn't initialized or file can't be opened
        }

        // Read the station ID (4 bytes) and store it as UInt32Bytes
        uint8_t buffer[BASE32_BYTES_LENGTH];
        if (file.read(buffer, BASE32_BYTES_LENGTH) != BASE32_BYTES_LENGTH) {
            file.close();
            return networkData;  // Return default values if station ID cannot be read
        }
        UInt32Bytes stationId(buffer);  // Store station ID as UInt32Bytes
        networkData = NetworkData(stationId);

        // Helper lambda to read a string with length prefix (without null termination)
        auto readString = [&](StringBytes &stringObj) {
            uint8_t length;
            if (file.read(&length, 1) != 1) {
                return false;  // Return false if length byte couldn't be read
            }
            uint8_t strBuffer[length];  // Buffer for exact length
            if (file.read(strBuffer, length) != length) {
                return false;  // Return false if string data couldn't be read
            }
            stringObj = StringBytes(length, strBuffer);  // Create StringBytes with length
            return true;
        };

        // Read IP
        StringBytes ip;
        if (!readString(ip)) {
            file.close();
            return networkData;  // Return default values if IP cannot be read
        }

        // Read Port
        StringBytes port;
        if (!readString(port)) {
            file.close();
            return networkData;  // Return default values if Port cannot be read
        }

        // Read SSID
        StringBytes ssid;
        if (!readString(ssid)) {
            file.close();
            return networkData;  // Return default values if SSID cannot be read
        }

        // Read Password
        StringBytes password;
        if (!readString(password)) {
            file.close();
            return networkData;  // Return default values if Password cannot be read
        }

        file.close();  // Close the file
        return NetworkData(stationId, ip, port, ssid, password);  // Return populated NetworkData object
    }

private:
    // Pins and driver for SoftSPI
    const uint8_t csPin, misoPin, mosiPin, sckPin;
    SoftSpiDriver<uint8_t, uint8_t, uint8_t> *softSpi;

    SdFat sd;
    SdFile file;
    bool sdInitialized;  // Stores the result of SD initialization
};

#endif  // SD_HANDLER_H
