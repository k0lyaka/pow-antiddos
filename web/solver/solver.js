const createHash = require("create-hash");

const generateNonce = (buf) => {
  let offset = 0;
  const timestamp = (Date.now() / 1000) >>> 0;

  buf[offset++] = (timestamp >> 24) & 0xff;
  buf[offset++] = (timestamp >> 16) & 0xff;
  buf[offset++] = (timestamp >> 8) & 0xff;
  buf[offset++] = timestamp & 0xff;

  const bits = offset + (((buf.length - offset) / 4) | 0) * 4;

  for (; offset < bits; offset++) buf[offset] = (Math.random() * 256) >>> 0;
};

const check = (hash, difficulty) => {
  let i = 0;
  let offset = 0;

  for (; i <= difficulty - 8; i += 8, offset++)
    if (hash[offset] !== 0) return false;

  const mask = 0xff << (8 + i - difficulty);
  return (hash[offset] & mask) === 0;
};

const solver = async (prefix, difficulty) => {
  const nonce = Buffer.alloc(16);

  while (true) {
    generateNonce(nonce);

    const hash = createHash("sha256")
      .update(prefix, "utf-8")
      .update(nonce)
      .digest();

    if (check(hash, difficulty)) return nonce.toString("hex");
  }
};

module.exports = solver;
