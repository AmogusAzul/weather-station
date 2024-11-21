#include "communication-bytes-types.h"

// Constructor HeaderBytes
Packet::Packet(const HeaderBytes& headerData) :
Bytes(headerData.getLength(), headerData.getBytes()),
header(headerData),
format(header.getField(HeaderField::VERSION)) {}

// Bulk Constructor
Packet::Packet(const Bytes& data) :
Bytes(data),
header(data.getBytes()),
format(header.getField(HeaderField::COUNT)) {}



// Set a new HeaderBytes without resizing the whole byte array
void Packet::setHeader(const HeaderBytes& newHeader) {
    header = newHeader;
    // Overwrite the header bytes in the beginning without resizing
    for (size_t i = 0; i < header.getLength(); ++i) {
        bytes[i] = header.getBytes()[i];
    }
}



void Packet::append(const Bytes& byteData, DataType type) {
    format.append(static_cast<uint8_t>(type));
    format.append(byteData.getLength());
    Bytes::append(byteData);
}