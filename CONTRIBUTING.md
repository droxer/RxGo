# Contributing to RxGo

We welcome contributions! Please follow these steps:

## How to Contribute

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Add** tests for new functionality
4. **Ensure** all tests pass (`make test`)
5. **Commit** your changes (`git commit -m 'Add amazing feature'`)
6. **Push** to the branch (`git push origin feature/amazing-feature`)
7. **Open** a pull request

## Development Setup

```bash
# Clone the repository
git clone https://github.com/droxer/RxGo.git
cd RxGo

# Install dependencies
make deps

# Run tests
make test

# Run benchmarks
make bench
```

## Development Guidelines

### Code Style
- Follow Go conventions and idioms
- Use meaningful variable and function names
- Add appropriate documentation for public APIs
- Ensure proper error handling

### Testing
- Write unit tests for all new functionality
- Ensure tests cover edge cases
- Use the table-driven test pattern when appropriate
- Add benchmarks for performance-critical code

### Pull Request Process
1. Ensure your branch is up to date with the main branch
2. Write clear, descriptive commit messages
3. Include tests for any new functionality
4. Update documentation if necessary
5. Ensure all CI checks pass

### Issue Reporting
- Use the GitHub issue tracker for bugs and feature requests
- Provide a clear description of the issue
- Include steps to reproduce for bugs
- Add relevant labels when possible

## Questions?

Feel free to open an issue for any questions about contributing to RxGo.