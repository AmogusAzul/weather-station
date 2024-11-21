#include "bytes-types.h"

// Constructor method (that takes a byte array)
Bytes::Bytes(const size_t length, const uint8_t* data)
{
    bytes.resize(length);
    if (data != nullptr) {
        setBytes(length, data);
    }
}

Bytes::Bytes(const std::initializer_list<Bytes> byteList)
{
        size_t length = 0;

        for (const auto& byteArray : byteList) {
            length += byteArray.getLength();
        }
        
        bytes.resize(size);
        for (const auto& byteArray : byteList) {
            bytes.insert(bytes.end(), byteArray.getBytes(), byteArray.getBytes() + byteArray.getLength());
        }
}

size_t Bytes::getLength() const
{
    return size;
}

void Bytes::setBytes(size_t length, const uint8_t* data)
{

    // matches the vector's length to array's 
    if (length != size) {
        bytes.resize(length);
        size = length;
    }

    // coping data to the vector
    for (int i = 0; i < length; i++) {
        bytes[i] = data[i];
    }

}

const uint8_t* Bytes::getBytes() const
{
    return bytes.data();

}

void Bytes::setByteString(std::string str)
{
    // Copy the string itself into the bytes array
    std::copy(str.begin(), str.end(), bytes.begin());
}

std::string Bytes::getByteString() const
{
    std::string str(bytes.begin(), bytes.end());
    return str;
}

uint8_t& Bytes::operator[](size_t index) {
    return bytes[index]; // Return reference to the byte
}

const uint8_t& Bytes::operator[](size_t index) const {
    return bytes[index]; // Return const reference to the byte
}

std::vector<std::pair<int, int>> 
Bytes::bigEndianBitShifts() const
{
    std::vector<std::pair<int, int>> pairs;
    for (int i = 0; i < getLength(); ++i) {
        pairs.emplace_back(i, 8 * (getLength() - i - 1));
    }
    return pairs;
}

// Append a single byte to the bytes vector
void Bytes::append(const uint8_t byte) {
    bytes.push_back(byte);
    size = bytes.size();  // Update size after appending
}

// Append an array of bytes to the bytes vector
void Bytes::append(const uint8_t* data, const size_t length) {
    bytes.insert(bytes.end(), data, data + length);  // Insert multiple bytes
    size = bytes.size();  // Update size after appending
}
void Bytes::append(const Bytes bytes) {
    append(bytes.getBytes(), bytes.getLength());
}