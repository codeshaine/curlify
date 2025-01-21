# Contributing to Curlify

Welcome, and thank you for considering contributing to Curlify! I appreciate your support and effort in making this project better for everyone. Whether you're fixing a bug, improving the UI, or adding new features, your contributions are always welcome.

## Getting Started

### Prerequisites

To contribute to Curlify, ensure the following tools are installed on your system:

- **Go**: Ensure that Go (1.19 or newer) is installed. You can download it [here](https://golang.org/dl/).
- **Git**: For version control to clone the repository and manage your contributions.
- **Bubble Tea**: Since I'm using bubble tea framework for developing this application, having good understanding of bubble tea will be a necessary
- **Make**: (Optional) For running build/test commands easily.

### Setting Up the Project

1. **Fork and Clone the Repository**:

   - Navigate to the [Curlify GitHub repository](https://github.com/codeshaine/curlify).
   - Click on the "Fork" button to create your own copy of the repository.
   - Clone your forked repository:
     ```bash
     git clone https://github.com/codeshaine/curlify.git
     cd curlify
     ```

2. **Create a Branch**:
   Create a branch using your username to work on your changes:

   ```bash
   git checkout -b <your-username>/<feature-name>
   ```

3. **Install Dependencies**:
   Run the following command to install dependencies:

   ```bash
   go mod tidy
   ```

4. **Run the Project**:
   In the root of the folder

   ```bash
   make run
   or
   go run .
   ```

5. **Run Tests**:
   Tests are critical for maintaining the quality of the project. Once the `Makefile` is added, you can run:

   ```bash
   make test
   or
   go test ./...
   ```

## How to Contribute

### Reporting Bugs

If you encounter any bugs or issues:

1. **Search Existing Issues**: Before opening a new issue, check if it has already been reported.
2. **Create a New Issue**: If it’s a new bug, [open an issue](https://github.com/codeshaine/curlify/issues/new) and include:
   - A clear title
   - Steps to reproduce the bug
   - Expected and actual behavior
   - Screenshots or logs, if applicable

### Suggesting Features

We welcome new ideas! To propose a new feature:

1. **Search Existing Issues**: Ensure your idea hasn’t been suggested already.
2. **Open a Feature Request**: Use the [issues page](https://github.com/codeshaine/curlify/issues/new) and select the "Feature Request" template.
3. **Be Descriptive**: Include details like why this feature is useful and how it might work.

### Contributing Code

1. **Create a New Branch**:
   Use a descriptive name for your branch:

   ```bash
   git checkout -b <your-username>/<feature-name>
   ```

```

```

2. **Make Changes**:
   Add your code changes to the appropriate files. Follow the project's coding conventions and comment where necessary.

3. **Run Tests**:
   Ensure your changes pass all tests:

   ```bash
   go test ./...
   ```

4. **Commit Your Changes**:
   Write clear and concise commit messages:

   ```bash
   git add .
   git commit -m "Feat: Description of your changes"
   ```

5. **Push Your Changes**:
   Push your changes to your forked repository:

   ```bash
   git push origin <your-username>/<feature-name>
   ```

6. **Create a Pull Request (PR)**:
   - Go to your forked repository on GitHub.
   - Click on "Compare & pull request."
   - Provide a clear description of your changes and link any related issues.

### Code Style

Please ensure that your code follows these guidelines:

- Use `go fmt` to format your code.
- Keep your code clean and readable.
- Write meaningful comments where necessary.

### Writing Tests

When adding new features or fixing bugs, include appropriate unit tests:

- Place tests in the corresponding `_test.go` file.
- Use descriptive test names and cover edge cases.

## Code of Conduct

By contributing to this project, you agree to abide by the rules outlined in our [Code of Conduct](https://github.com/codeshaine/curlify/blob/main/CODE_OF_CONDUCT.md). Please take a moment to read it.

## Community Guidelines

- Be kind and respectful to other contributors.
- Use inclusive language in code, comments, and discussions.
- Follow GitHub's [Community Guidelines](https://docs.github.com/en/github/site-policy/github-community-guidelines).

## Need Help?

If you have any questions or need assistance, feel free to:

- Open a [GitHub discussion](https://github.com/codeshaine/curlify/discussions).
- Tag `@codeshaine` in your issue or PR for urgent queries.

Thank you for contributing to Curlify! Together, we can make this project amazing.
