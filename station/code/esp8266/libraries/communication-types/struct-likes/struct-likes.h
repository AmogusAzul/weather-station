#ifndef STRUCT_LIKES_H
#define STRUCT_LIKES_H

#include "../../common/types/base32-bytes-types/base32-bytes-types.h"
#include "../../common/types/string-bytes-types/string-bytes-types.h"

class NetworkData {
public:
    // Constructor with default values
    NetworkData(const UInt32Bytes &stationId = UInt32Bytes(), 
            const StringBytes &ip = StringBytes(),
            const UInt32Bytes &port = UInt32Bytes(),
            const StringBytes &ssid = StringBytes(),
            const StringBytes &password = StringBytes());

    // Getter methods
    UInt32Bytes getStationId() const;
    StringBytes getIp() const;
    UInt32Bytes getPort() const;
    StringBytes getSsid() const;
    StringBytes getPassword() const;

private:
    UInt32Bytes stationId;  // Store station ID as UInt32Bytes
    StringBytes ip;
    UInt32Bytes port;
    StringBytes ssid;
    StringBytes password;
};

#endif  // STRUCT_LIKES_H