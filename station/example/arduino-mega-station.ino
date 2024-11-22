#include <Arduino.h>

#define SERIAL_BAUD 115200

// Hardcoded parameters for the sine wave
const float amplitude = 100.0;   // Amplitude of the sine wave
const float frequency = 10.0;   // Frequency of the sine wave
const float noiseLevel = 0.2;  // Level of noise to add to the sine wave

// Calculate the period (number of steps for one complete cycle)
const int period = (int)(2 * PI / frequency);  // Number of steps for the wave to repeat

// Function to generate a sine wave value with noise and control step
uint32_t generateSineWaveWithNoise(int &step) {
    // Calculate the angle using the step and frequency
    float angle = step * frequency;

    // Calculate the sine wave value at this step
    float sineValue = amplitude * sin(angle);

    // Add noise to the sine wave
    float noise = (random(-100, 100) / 100.0) * noiseLevel;
    float sineWithNoise = sineValue + noise;

    // Clamp sineWithNoise to ensure it stays within the range [-1, +1]
    if (sineWithNoise > 1.0) sineWithNoise = 1.0;
    if (sineWithNoise < -1.0) sineWithNoise = -1.0;

    // Normalize sine value to range from 0 to a uint32_t maximum
    uint32_t normalizedValue = (uint32_t)((sineWithNoise + 1) * (INT32_MAX / 2));

    // Update step value and wrap around after one complete cycle
    step = (step + 1) % period;

    return normalizedValue;
}

// Function to send an `int32_t` in big-endian format
void sendInt32BigEndian(uint32_t value) {
  byte buffer[4];
  
  // Convert to big-endian (most significant byte first)
  buffer[0] = (value >> 24) & 0xFF;
  buffer[1] = (value >> 16) & 0xFF;
  buffer[2] = (value >> 8) & 0xFF;
  buffer[3] = value & 0xFF;

  // Send each byte to the ESP8266
  Serial3.write(buffer, sizeof(buffer));
}

void setup() {
  // Initialize the serial port
  Serial.begin(SERIAL_BAUD);
  Serial3.begin(SERIAL_BAUD);
  delay(5000);
}

void loop() {
  // Generate and send a random number every 5 seconds
  int step = 0;
  uint32_t randomNumber = generateSineWaveWithNoise(step);
  
  sendInt32BigEndian(randomNumber);  // Custom function to send the number

  // Forward data from Serial3 to Serial0
  if (Serial3.available() > 0) {  // Check if data is available on Serial3
    while (Serial3.available() > 0) {
      char incomingByte = Serial3.read();  // Read one byte from Serial3
      Serial.write(incomingByte);          // Send it to Serial0 (Serial Monitor)
    }
  }

  // Wait for 5 seconds before sending the next number
  delay(5000);
}
