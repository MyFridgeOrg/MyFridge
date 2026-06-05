# Commit Message Convention

This project follows the **Semantic Commit Messages** convention as described by [joshbuchea](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716).

## Format

```
<type>: <subject>
```

## Types

| Type       | Description                                                |
|------------|------------------------------------------------------------|
| `feat`     | A new feature                                              |
| `fix`      | A bug fix                                                  |
| `docs`     | Documentation only changes                                 |
| `style`    | Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc.) |
| `refactor` | A code change that neither fixes a bug nor adds a feature  |
| `perf`     | A code change that improves performance                    |
| `test`     | Adding missing tests or correcting existing tests          |
| `chore`    | Changes to the build process or auxiliary tools and libraries such as documentation generation |

## Examples

```
feat: add hat wobble
fix: resolve off-by-one error in loop
docs: update README with new API endpoints
style: convert tabs to spaces
refactor: extract shared logic into helper function
perf: reduce bundle size by lazy-loading modules
test: add unit tests for user service
chore: update CI pipeline configuration
```

## Rules

- Use the **imperative, present tense** ("add", not "added" or "adds").
- **Do not** capitalize the subject line.
- **Do not** end the subject line with a period.
- Keep the subject line concise (≤ 72 characters).

