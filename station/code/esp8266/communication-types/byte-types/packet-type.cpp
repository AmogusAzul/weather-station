#include "communication-bytes-types.h"

// Constructor using initializer list that includes both the HeaderBytes and other data
Packet::Packet(const HeaderBytes& headerData, const std::initializer_list<Bytes> byteList) : 
Bytes::Bytes(headerData.getLength(), headerData.getBytes()), header(headerData) {
    // Append the rest of the data after the header
    for (const Bytes& byteChunk : byteList) {
        append(byteChunk.getBytes(), byteChunk.getLength());
    }
}

// Set a new HeaderBytes without resizing the whole byte array
void Packet::setHeader(const HeaderBytes& newHeader) {
    header = newHeader;
    // Overwrite the header bytes in the beginning without resizing
    for (size_t i = 0; i < header.getLength(); ++i) {
        bytes[i] = header.getBytes()[i];
    }
}