#include "struct-likes.h"

// Constructor with default values
NetworkData::NetworkData(const UInt32Bytes &stationId, const StringBytes &ip, const StringBytes &port,
                         const StringBytes &ssid, const StringBytes &password)
    : stationId(stationId), ip(ip), port(port), ssid(ssid), password(password) {}

// Getter methods
UInt32Bytes NetworkData::getStationId() const {
    return stationId;
}

StringBytes NetworkData::getIp() const {
    return ip;
}

StringBytes NetworkData::getPort() const {
    return port;
}

StringBytes NetworkData::getSsid() const {
    return ssid;
}

StringBytes NetworkData::getPassword() const {
    return password;
}