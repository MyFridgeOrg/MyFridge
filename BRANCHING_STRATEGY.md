# Branching Strategy

This project uses an **Epic / Feature** branching flow.

## Branch Types

| Branch              | Purpose                                      |
|---------------------|----------------------------------------------|
| `main`              | Production-ready code. Releases are cut from here. |
| `epic/{epic-name}`  | Groups related features that together deliver a larger body of work. |
| `feature/{feature-name}` | A single, self-contained piece of work within an epic. |

## Flow

```
main
 └── epic/user-management
      ├── feature/login
      ├── feature/registration
      └── feature/password-reset
```

1. **Create an epic branch** from `main`:
   ```bash
   git checkout main
   git checkout -b epic/user-management
   ```

2. **Create feature branches** from the epic:
   ```bash
   git checkout epic/user-management
   git checkout -b feature/login
   ```

3. **Merge features into the epic** — once a feature is complete and reviewed, merge it back into its epic branch.

4. **Merge the epic into `main`** — when all features in the epic are complete and the epic is fully tested, merge the epic branch into `main`.

5. **Release from `main`** — production releases are always cut from the `main` branch.

## Summary

```
feature/* ──► epic/* ──► main ──► release
```

