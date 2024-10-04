#ifndef BYTES_TYPES_H
#define BYTES_TYPES_H

#include <vector>
#include <cstddef>  // For size_t
#include <cstdint>  // For uint8_t, int32_t
#include <string>   // For string


class Bytes {

protected:

    // Declares a vector for the data
    std::vector<uint8_t> bytes;

    // Size of the byte array
    size_t size;

public:
    // Constructor (just a declaration here)
    Bytes(const size_t length, const uint8_t* data = nullptr);

    // Get the length of the byte array
    size_t getLength() const;

    // Set data from a byte array
    void setBytes(const size_t length, const uint8_t* data);

    // Get a pointer to the byte array
    const uint8_t* getBytes() const;

    void setByteString(std::string str);
    std::string getByteString() const;

    // Method to access individual bytes
    uint8_t& operator[](const size_t index);
    const uint8_t& operator[](const size_t index) const;

    // Function for getting corresponding bitshifts to bytes' indexes
    std::vector<std::pair<int, int>>bigEndianBitShifts() const;
};

#endif // BYTES_TYPES_H