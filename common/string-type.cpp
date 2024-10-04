#include "headers/types/string-bytes-types/string-bytes-types.h"

#include <algorithm> // For std::copy, std::min, std::max

// Making sure that the string is under the STRING_MAX_LEVEL constant
StringBytes::StringBytes(std::string str) :
Bytes::Bytes(
    (str.length() > STRING_MAX_LENGTH ? 
    STRING_MAX_LENGTH : str.length())
    + STRING_LENGTH_BYTES_LENGTH
    )
{
    setString(str);
}

void StringBytes::setBytes(size_t length, uint8_t *data)
{

    size = std::min(length, STRING_MAX_LENGTH + STRING_LENGTH_BYTES_LENGTH);
    size = std::max(size, STRING_LENGTH_BYTES_LENGTH);

    Bytes::setBytes(size, data);

}

void StringBytes::setString(const std::string str)
{

    // Cap the str to STRING_MAX_LENGTH
    std::string cappedStr = str.length() > STRING_MAX_LENGTH ? 
    str.substr(0, STRING_MAX_LENGTH) : str;

    // Update size to cappedStr's length 
    size = cappedStr.length() + STRING_LENGTH_BYTES_LENGTH;
    bytes.resize(size);

    // iterating over every byte of the size masking it with "0xFF"
    for (const auto& pair : bigEndianBitShifts()) {
        bytes[pair.first] = (size >> pair.second) & 0xFF;
    }

    // Copy the string itself into the rest of the bytes array
    std::copy(cappedStr.begin(), cappedStr.end(), bytes.begin() + STRING_LENGTH_BYTES_LENGTH);

}

std::string StringBytes::getString() const
{
    // Extract the string content from the remaining bytes
    std::string str(
        bytes.begin() + STRING_LENGTH_BYTES_LENGTH, 
        bytes.begin() + STRING_LENGTH_BYTES_LENGTH + size);

    return str;
}

size_t StringBytes::getStringLength() const
{
    return getLength() - STRING_LENGTH_BYTES_LENGTH;
}
