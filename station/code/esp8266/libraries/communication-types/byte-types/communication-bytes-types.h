#ifndef COMMUNICATION_BYTES_TYPES
#define COMMUNICATION_BYTES_TYPES

#include "../../common/types/bytes-types.h"
#include <string>

class TokenBytes : public Bytes {
public:

    TokenBytes(const uint8_t* data = (const uint8_t*)nullptr);
    TokenBytes(const std::string str);

};

enum class HeaderField : size_t {
    VERSION = 0,
    TYPE = 1,
    SPECIFIC_TYPE = 2,

    COUNT
};

class HeaderBytes : public Bytes {
public:

    HeaderBytes(const uint8_t* data);
    HeaderBytes(const std::string str);
    HeaderBytes(const std::initializer_list<uint8_t> fields);

    // Unified getter and setter
    void setField(HeaderField field, uint8_t value);
    uint8_t getField(HeaderField field) const;

};

// Get the count of enum values
static constexpr size_t getFieldEnum(HeaderField field) {
    return static_cast<size_t>(field);
}

enum class DataType : uint8_t {

    BYTES = 0,
    INT32 = 1,
    UINT32 = 2,
    FLOAT32 = 3,
    STRING = 4,
    TOKEN = 5,
};

class Packet : public Bytes {
private:
    HeaderBytes header;  // All packets start with a HeaderBytes object

    Bytes format;

public:
    // Constructor using initializer list that includes both the HeaderBytes and other data
    Packet(const HeaderBytes& headerData);

    Packet(const Bytes& data);

    // Get a reference to the HeaderBytes
    const HeaderBytes& getHeader() const {
        return header;
    }

    void setHeader(const HeaderBytes& newHeader);

    void append(const Bytes& bytes, DataType type);

    Bytes getFormat() const {
        return format;
    };

};

#endif // COMMUNICATION_BYTES_TYPES