#include "communication-bytes-types.h"

HeaderBytes::HeaderBytes(const uint8_t* data) :
Bytes::Bytes(getFieldEnum(HeaderField::COUNT), data)
{
}

HeaderBytes::HeaderBytes(const std::string str) :
Bytes::Bytes(getFieldEnum(HeaderField::COUNT))
{
    setByteString(str);
}

HeaderBytes::HeaderBytes(const std::initializer_list<uint8_t> fields) :
Bytes::Bytes(getFieldEnum(HeaderField::COUNT))
{
    // Fill the fields with provided values and default to 0 for missing fields
    size_t i = 0;
    for (auto field : fields) {
        if (i < getFieldEnum(HeaderField::COUNT)) {
            bytes[i++] = field;
        }
    }

    // Set any remaining fields to 0
    while (i < getFieldEnum(HeaderField::COUNT)) {
        bytes[i++] = 0;
    }
}

void HeaderBytes::setField(HeaderField field, uint8_t value)
{
    bytes[getFieldEnum(field)] = value;
}

uint8_t HeaderBytes::getField(HeaderField field) const
{
    return bytes[getFieldEnum(field)];
}
