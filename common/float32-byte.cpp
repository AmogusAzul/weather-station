#include "headers/types/base32-bytes-types/base32-bytes-types.h"
#include <cstring>

Float32Bytes::Float32Bytes(const float number) :
Base32Bytes::Base32Bytes()
{
    setNumber(number);
};

Float32Bytes::Float32Bytes(const uint8_t* data) : 
Base32Bytes::Base32Bytes(data)
{
}

void Float32Bytes::setNumber(const float number)
{
    // Reinterpret the float as uint32_t for bit manipulation
    uint32_t dummyInt;
    std::memcpy(&dummyInt, &number, sizeof(float));  // Alternative to reinterpret_cast
    
    // iterating over every byte of the int32_t and masking it with "0xFF"
    for (const auto& pair : bigEndianBitShifts()) {
        bytes[pair.first] = (dummyInt >> pair.second) & 0xFF;
    }
}

float Float32Bytes::getNumber() const
{
    uint32_t dummyInt;

    //populate dummyInt with the bytes' data
    // OR-ing every shifted value of the vector
    for (const auto& pair : bigEndianBitShifts()) {

        dummyInt |= (static_cast<int32_t>(bytes[pair.first]) << pair.second);
        
    }

    // Reinterpret the bit representation as a float
    float number;
    std::memcpy(&number, &dummyInt, sizeof(float));


    return number;
}
