#include "libraries/common/types/base32-bytes-types/base32-bytes-types.h"
#include "libraries/communication-types/byte-types/communication-bytes-types.h"
#include "libraries/communication-types/struct-likes/struct-likes.h"

#include "libraries/common/types/bytes-type.cpp"

#include "libraries/common/types/base32-bytes-types/base32-bytes-type.cpp"
#include "libraries/common/types/base32-bytes-types/uint32-type.cpp"

#include "libraries/common/types/string-bytes-types/string-type.cpp"

#include "libraries/communication-types/byte-types/token-type.cpp"
#include "libraries/communication-types/byte-types/header-type.cpp"
#include "libraries/communication-types/byte-types/packet-type.cpp"

#include "libraries/communication-types/struct-likes/network-data.cpp"



#include <vector>


#include "libraries/sd-handler/sd-handler.hpp"
#include "libraries/network-client/network-client.hpp"

#include <SPI.h>
#include "SdFat.h"
#include <Arduino.h>
#include <cstdint>
#include <string>

#define MISO_PIN 12
#define MOSI_PIN 4
#define SCK_PIN 5
#define CHIP_SELECT_PIN 16

#define REQUEST_VERSION ((uint8_t)(1))
#define REQUEST_TYPE ((uint8_t)(1))
#define SPECIFIC_REQUEST_TYPE ((uint8_t)(0))

const std::string ip = "192.168.110.153";
const std::string ssid = "FLIA ANGARITA";
const std::string password = "SEBDARRO";

const uint8_t stationId = 1;


const byte expectedData[] = {
  static_cast<byte>(DataType::UINT32), BASE32_BYTES_LENGTH, 
  static_cast<byte>(DataType::UINT32), BASE32_BYTES_LENGTH
};
// Define dataSize as the count of elements in expectedData
constexpr size_t dataSize = sizeof(expectedData) / sizeof(expectedData[0]);

size_t getTotalSize(const byte* array, size_t length){
  size_t result;
  for (int i=1; i<length;i+=2){
    result+=array[i];
  }
  return result;
}

const size_t fullLength = getTotalSize(expectedData, dataSize);


// Define SDHandler object using template arguments and initialize with CHIP_SELECT_PIN
SDHandler<MISO_PIN, MOSI_PIN, SCK_PIN> sd(CHIP_SELECT_PIN);

int tryTest() {
  Serial.begin(115200);
  Serial.println("\npear");
  return 3;
}

int test = tryTest();

Bytes testBytes(40);


NetworkClient client(netData);

TokenBytes token;

void setup() {
  Serial.begin(115200);
  
  // Initialize SD card and get configuration
  netData = getNetData();
  token = sd.readToken();
}

void loop() {
  // Wait until we receive expected data
  while (Serial.available() < fullLength) {
    yield();
  }

  // Read Serial data in smaller chunks
  byte buffer[fullLength];
  for (size_t i = 0; i < fullLength; ++i) {
    buffer[i] = Serial.read();
    yield();  // Yield between reads to avoid long blocking
  }

  // Separate the buffer based on expectedData array
  size_t offset = 0;
  UInt32Bytes randomNum(buffer + offset);
  offset += BASE32_BYTES_LENGTH;
  UInt32Bytes timestamp(buffer + offset);
  offset += BASE32_BYTES_LENGTH;

  HeaderBytes header({REQUEST_VERSION, REQUEST_TYPE, SPECIFIC_REQUEST_TYPE});
  Packet packet(header);

  packet.append(netData.getStationId(), DataType::UINT32);
  packet.append(token, DataType::TOKEN);
  packet.append(timestamp, DataType::UINT32);
  packet.append(randomNum, DataType::UINT32);

  auto [isConnected, justRecovered] = client.canSend();

  if (isConnected) {
    Packet responsePacket(nullptr);
    client.send(packet, responsePacket);

    HeaderBytes responseHeader(responsePacket.getHeader());

    if (responseHeader.getField(HeaderField::VERSION) != REQUEST_VERSION) {
      client.forceBackup();
      return;
    }

    if (responseHeader.getField(HeaderField::TYPE) != 1) {
      switch (responseHeader.getField(HeaderField::SPECIFIC_TYPE)) {
        case 1:
        case 4:
        case 5:
        case 6:
        case 7:
          client.forceBackup();
          return;
      }
    }

    token = TokenBytes(responsePacket.getBytes() + static_cast<size_t>(HeaderField::COUNT));
    sd.swapToken(token);
  }

  // Backup to SD if not connected
  if (!isConnected && sd.isSDInitialized()) {
    sd.writeBackupInfo(packet);
  }
}
