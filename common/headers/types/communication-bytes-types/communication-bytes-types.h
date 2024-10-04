#ifndef COMMUNICATION_BYTES_TYPES
#define COMMUNICATION_BYTES_TYPES

#include "../bytes-types.h"
#include <string>

class TokenBytes : public Bytes {

    TokenBytes(uint8_t* data = nullptr);
    TokenBytes(std::string str);

};

enum class HeaderField : size_t {
    VERSION = 0,
    TYPE = 1,
    DIRECTION = 2,

    COUNT
};

class HeaderBytes : public Bytes {

    HeaderBytes(uint8_t* data);
    HeaderBytes(std::string str);
    HeaderBytes(std::initializer_list<uint8_t> fields);

    // Unified getter and setter
    void setField(HeaderField field, uint8_t value);
    uint8_t getField(HeaderField field) const;

    // Get the count of enum values
    static constexpr size_t getFieldEnum(HeaderField field) {
        return static_cast<size_t>(field);
    }

};

#endif // COMMUNICATION_BYTES_TYPES