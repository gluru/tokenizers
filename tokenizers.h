#include <stdbool.h>
#include <stdint.h>

void *from_bytes(const uint8_t *config, uint32_t len);

void *from_bytes_with_truncation(const uint8_t *config, uint32_t len, uint32_t max_len, uint8_t direction);

void *from_file(const char *config);

void free_tokenizer(void *ptr);

uint32_t *encode(void *ptr, const char *message, uint32_t *len, bool add_special_tokens);

char *decode(void *ptr, const uint32_t *ids, uint32_t len, bool skip_special_tokens);

uint32_t vocab_size(void *ptr);
