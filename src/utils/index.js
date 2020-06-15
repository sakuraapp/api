exports.padLeft = function (str, len = 4) {
    return Array(len - String(str).length + 1).join('0') + str
}
