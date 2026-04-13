# Hangman Arena

**hangman-arena** is a high-performance **Go CLI tool** designed to benchmark and evaluate the linguistic reasoning capabilities of **Large Language Models (LLMs)**. By forcing models to play the game of Hangman, the tool strips away conversational fluff and tests a model's ability to handle **state tracking**, **negative constraints**, and **conditional probability** at a character level.

---

### **Key Features**

- **Multi-Agent Benchmarking:** Supports diverse "Player" types, including **Random Character Guessers** and **Ollama-powered LLMs** (e.g., Gemma 4 2b, 31b).
- **Stateless Execution:** Each guess is treated as a fresh request to prevent "Context Bloat," ensuring the model relies on the current board state rather than past conversational noise.
- **JSON-Enforced Reasoning:** Utilizes Ollama's format constraints to force models to provide both a **move** and a **logical justification**, allowing for deep analysis of the model's "internal thought process."
- **Performance Analytics:** Automatically generates detailed `.jsonl` logs that track win rates, total attempts, and specific failure modes like "Length Violations" or "Illegal Moves."

---

### **Core Benchmarking Findings**

Initial benchmarks using **Gemma 4 2b** revealed critical insights into small-model behavior:

- **Linguistic Hard-Coding:** Smaller models often default to a rigid, global frequency list (e.g., guessing 'e', 't', 'a', 's', 'i', 'o') rather than adapting to the specific word length.
- **Length Constraint Failures:** Models frequently hallucinate solutions that do not match the required character count, such as suggesting the 6-letter word "column" for an 8-letter pattern.
- **Negative Constraint Struggles:** While models can identify the current pattern, they occasionally suggest letters in their "reasoning" that have already been marked as incorrect.
- **Pattern Matching Success:** Models show strong proficiency in identifying common English suffixes (like `-ound`) once enough structural consonants are revealed.

---

### **Why it Matters**

Most LLM benchmarks focus on high-level prose or coding. **hangman-arena** focuses on **Reverse Tokenization**—the difficult task of breaking down concept-level tokens (like "cow") into individual character components. This project highlights the "Parameter Ceiling" where models transition from stochastic pattern matching to genuine linguistic intuition.