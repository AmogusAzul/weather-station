#ifndef COMMUNICATION_BYTES_TYPES
#define COMMUNICATION_BYTES_TYPES

#include "../../../common/types/bytes-types.h"
#include <string>

class TokenBytes : public Bytes {

    TokenBytes(const uint8_t* data = (const uint8_t*)nullptr);
    TokenBytes(const std::string str);

};

enum class HeaderField : size_t {
    VERSION = 0,
    TYPE = 1,
    DIRECTION = 2,

    COUNT
};

class HeaderBytes : public Bytes {

    HeaderBytes(const uint8_t* data);
    HeaderBytes(const std::string str);
    HeaderBytes(const std::initializer_list<uint8_t> fields);

    // Unified getter and setter
    void setField(HeaderField field, uint8_t value);
    uint8_t getField(HeaderField field) const;

    // Get the count of enum values
    static constexpr size_t getFieldEnum(HeaderField field) {
        return static_cast<size_t>(field);
    }

};

class Packet : public Bytes {
private:
    HeaderBytes header;  // All packets start with a HeaderBytes object

public:
    // Constructor using initializer list that includes both the HeaderBytes and other data
    Packet(const HeaderBytes& headerData, const std::initializer_list<Bytes> byteList);

    // Get a reference to the HeaderBytes
    const HeaderBytes& getHeader() const {
        return header;
    }

    void Packet::setHeader(const HeaderBytes& newHeader);

};

#endif // COMMUNICATION_BYTES_TYPES