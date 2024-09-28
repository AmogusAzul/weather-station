#ifndef COMMUNICATION_BYTES_TYPES
#define COMMUNICATION_BYTES_TYPES

#include "../bytes-types.h"
#include <string>

class TokenBytes : public Bytes {

    TokenBytes(uint8_t* data = nullptr);
    TokenBytes(std::string str);

};

#endif // COMMUNICATION_BYTES_TYPES