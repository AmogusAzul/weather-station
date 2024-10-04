#include "headers/types/base32-bytes-types/base32-bytes-types.h"

Int32Bytes::Int32Bytes(int32_t number) : 
Base32Bytes::Base32Bytes()
{
    setNumber(number);
}

Int32Bytes::Int32Bytes(const uint8_t* data) : 
Base32Bytes::Base32Bytes(data)
{
}

void Int32Bytes::setNumber(const int32_t number)
{
    // iterating over every byte of the int32_t masking it with "0xFF"
    for (const auto& pair : bigEndianBitShifts()) {
        bytes[pair.first] = (number >> pair.second) & 0xFF;
    }
}

int32_t Int32Bytes::getNumber() const
{
    int32_t number;

    // OR-ing every shifted value of the vector
    for (const auto& pair : bigEndianBitShifts()) {

        number |= (static_cast<int32_t>(bytes[pair.first]) << pair.second);
    }

    return number;
}
