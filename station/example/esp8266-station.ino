#include <SD.h>


#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <time.h>

const char* ssid = "FLIA ANGARITA";
const char* password = "SEBDARRO";

// TCP server details
const char* tcp_server_ip = "192.168.110.153";  // Replace with your server's IP
const uint16_t tcp_server_port = 8080;          // Replace with your server's port

WiFiClient client;

// Variables for the station request
int32_t stationID = 1;  // Example station ID

const int tokenLength = 6;                                   // Length of the token expected by the server
byte token[tokenLength] = { '6', '9', '6', '9', '6', '9' };  // ASCII values for '696969'

void setup() {
  Serial.begin(115200);  // For communication with the Arduino Mega

  // Connect to WiFi
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
  }
}

void loop() {
  // Wait for a new message from Arduino
  waitForMessage();

  // After receiving data, send it to the server
}

void waitForMessage() {
  // Wait until we receive at least 4 bytes from the serial buffer
  while (Serial.available() < 4) {
    delay(100);  // Add a small delay to prevent high CPU usage
  }

  // Read the 4 bytes from the serial buffer
  byte buffer[4];
  Serial.readBytes(buffer, 4);  // Read 4 bytes

  // Assemble the bytes into a 32-bit integer (big-endian)
  int32_t randomNum = ((int32_t)buffer[0] << 24) | ((int32_t)buffer[1] << 16) | ((int32_t)buffer[2] << 8) | buffer[3];

  // Send the station request with the random number to the server
  sendStationRequest(randomNum);
}

// Function to merge and send the station request as one large packet
void sendStationRequest(int32_t randomNum) {
  if (client.connect(tcp_server_ip, tcp_server_port)) {

    // Create a buffer to hold the entire message
    const int packetSize = 1 + 1 + 1 + 4 + tokenLength + 4 + 4;  // Version (1B), Request Type (1B), Station ID (4B), Token (12B), Timestamp (4B), Random Num (4B)
    byte packet[packetSize];
    int index = 0;

    // 1. Add Version (1 byte)
    packet[index++] = 1;

    // 2. Add Request Type (1 byte)
    packet[index++] = 1;
    
    // 2.1. Add Specific Request Type (1 byte)
    packet[index++] = 1;

    // 3. Add Station ID (4 bytes, big-endian)
    packet[index++] = (stationID >> 24) & 0xFF;
    packet[index++] = (stationID >> 16) & 0xFF;
    packet[index++] = (stationID >> 8) & 0xFF;
    packet[index++] = stationID & 0xFF;

    // 4. Add token
    for (int i = 0; i < tokenLength; i++) {
      packet[index++] = token[i];
    }

    // 5. Add current Unix timestamp (4 bytes, big-endian)
    int32_t timestamp = (int32_t)time(nullptr);  // Get current Unix time
    packet[index++] = (timestamp >> 24) & 0xFF;
    packet[index++] = (timestamp >> 16) & 0xFF;
    packet[index++] = (timestamp >> 8) & 0xFF;
    packet[index++] = timestamp & 0xFF;

    // 6. Add Random Number (4 bytes, big-endian)
    packet[index++] = (randomNum >> 24) & 0xFF;
    packet[index++] = (randomNum >> 16) & 0xFF;
    packet[index++] = (randomNum >> 8) & 0xFF;
    packet[index++] = randomNum & 0xFF;

    // Send the complete packet as a single message
    client.write(packet, packetSize);

    // Wait for the server to respond with the full buffer (version, type, specific type, token if ok)
    while (!client.available()) {
      delay(100);  // Wait for a response
    }

    // In the sendStationRequest function:
    byte responseBuffer[2 + 1 + tokenLength];  // 1B version + 1B response type + 1B specific type + token (if present)
    int responseSize = client.read(responseBuffer, sizeof(responseBuffer));

    // Parse the response
    byte responseVersion = responseBuffer[0];
    byte responseType = responseBuffer[1];
    byte specificType = responseBuffer[2];

    // Check the version and response type
    if (responseVersion != 1) {
      client.stop();
      return;
    }

    if (responseType == 2 && specificType == 2) {  // okType = 2 and specificType = 02 (only return token)
      // Handle success response (okType = 2 with token)
      memcpy(token, &responseBuffer[3], tokenLength);  // Copy the token
    } else if (responseType == 1 && specificType == 1) {
    } else {
      // Unexpected response type or specific type
    }

    // Close the connection
    client.stop();
  }
}

// Function to convert an int32_t into its 4-byte big-endian representation as a string
String int2BEStr(int32_t num) {
  char rawBytes[4];  // Array to store 4 raw bytes

  // Extract each byte in big-endian order and store it in rawBytes
  for (int i = 0; i < 4; i++) {
    rawBytes[3 - i] = (num >> (i * 8)) & 0xFF;  // Shift and mask to get each byte
  }

  return String(rawBytes);  // Convert the 4-byte array into a String
}
