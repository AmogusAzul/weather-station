#ifndef STRING_BYTES_TYPES_H
#define STRING_BYTES_TYPES_H

#include "../bytes-types.h"

#include <string>

// How many bytes will be used to store the length of the string
const size_t STRING_LENGTH_BYTES_LENGTH = 2;
// Consequent max length
const size_t STRING_MAX_LENGTH = (1UL << (8 * STRING_LENGTH_BYTES_LENGTH));

class StringBytes : public Bytes {
public:

    StringBytes(std::string str="");
    StringBytes(size_t length, uint8_t* data = nullptr);

    void setBytes(size_t length, uint8_t* data);

    void setString(std::string str);
    std::string getString() const;

    size_t getStringLength() const;

};

#endif // STRINGBYTES_TYPES_H