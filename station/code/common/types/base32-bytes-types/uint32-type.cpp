#include "base32-bytes-types.h"

UInt32Bytes::UInt32Bytes(uint32_t number) : 
Base32Bytes::Base32Bytes()
{
    setNumber(number);
}

UInt32Bytes::UInt32Bytes(const uint8_t* data) : 
Base32Bytes::Base32Bytes(data)
{
}

void UInt32Bytes::setNumber(const uint32_t number)
{
    // iterating over every byte of the int32_t and masking it with "0xFF"
    for (const auto& pair : bigEndianBitShifts()) {
        bytes[pair.first] = (number >> pair.second) & 0xFF;
    }
}

uint32_t UInt32Bytes::getNumber() const
{
    uint32_t number;

    // OR-ing every shifted value of the vector
    for (const auto& pair : bigEndianBitShifts()) {

        number |= (static_cast<uint32_t>(bytes[pair.first]) << pair.second);

    }

    return number;
}

