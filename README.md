# Go-Ratchet 🛡️

Double ratchet algorithm implementation.

# TODO

- Add docs for each function. Add comments for e.g. ratchetReceivingChain and ratchetSendingChain.
- Add tests.
- Reduce allocations count. For example, reuse slices for HKDF and encryption/decryption. Encrypt/Decrypt to array from stack.
- Create benchmarks to increase speed.
