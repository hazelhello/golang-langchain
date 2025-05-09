package chain

type CacheProvider interface {
    Get(key string) (string, bool)
    Set(key, value string)
}

type LLMClient interface {
    Call(prompt string) (string, error)
}