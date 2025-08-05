# Minimalistic Planning Poker - Coding Guidelines
Welcome to the team! This document provides all the necessary guidelines to contribute to this project. Our goal is to create a simple, fast, and reliable planning poker tool with minimal dependencies and a clean, maintainable codebase. Adhering to these principles is crucial.

## 1. Guiding Principles
Minimalism: Keep it simple. Every feature, dependency, and line of code should be essential. If it can be done simply, do it simply.

Readability: Write code for humans first, machines second. Use clear names for variables, functions, and files. Add comments only when the code's purpose isn't obvious.

No Frameworks (Frontend): We are not using any JavaScript frameworks (like React, Vue, or Angular). All frontend code will be written in vanilla JavaScript, HTML, and CSS. This reduces bloat and increases performance.

Zero-to-Minimal Dependencies: Only add a dependency if it's absolutely necessary and there isn't a straightforward way to implement the functionality within the project. Every dependency is a liability.

## 2. Technology Stack
Backend: Go

Frontend: Vanilla JavaScript (ES6+), HTML5, CSS3

API: WebSockets for real-time communication, with a simple REST API for session management if needed.

## 3. Backend Development (Go)
The backend is the core of our application. It handles game state, user connections, and business logic. We use Go for its performance, simplicity, and strong concurrency model.

### 3.1. Hexagonal Architecture (Ports & Adapters)
We follow the Hexagonal Architecture (also known as Ports and Adapters) to ensure a clean separation of concerns. This makes the application easier to test, maintain, and evolve. The project's directory structure is organized to reflect this architecture. 

### 3.2. Unit Testing
Writing unit tests for the backend is not optional; it is mandatory.

Test the Core: All business logic within /application must have 100% test coverage. Use mocks for the repository ports to test the services in isolation.

Test Adapters: Write tests for your adapters to ensure they correctly implement the port contracts.

File Naming: Test files must be named _test.go (e.g., session_service_test.go).

Keep Tests Simple: A test should focus on one thing only. Use descriptive test function names, like TestCreateSession_Success or TestAddUser_SessionIsFull.

## 4. Frontend Development
   The frontend should be clean, simple, and efficient.

### 4.1. Code Structure
Organize files logically. All frontend assets should be in a /frontend or /static directory served by the Go backend.

/frontend
├── /css
│   └── style.css
├── /js
│   ├── main.js         // Main entry point, event listeners
│   ├── api.js          // Handles WebSocket connection and messages
│   └── ui.js           // DOM manipulation functions
└── index.html

### 4.2. JavaScript (Vanilla)
Use Modules: Use ES6 modules (import/export) to keep your code organized.

State Management: Maintain the application state in a single object in main.js. When the state changes, call functions in ui.js to re-render the necessary parts of the DOM. Do not store state in the DOM (e.g., in data- attributes).

DOM Manipulation:

Create functions in ui.js dedicated to updating the DOM (e.g., renderUserList(users), showVotingResults(results)).

Avoid direct manipulation in main.js or api.js.

Use document.getElementById or document.querySelector to select elements. Cache frequently accessed elements in variables.

### 4.3. HTML & CSS
Semantic HTML: Use meaningful HTML5 tags (<main>, <section>, <nav>, etc.).

CSS: Do not write custom CSS, use Tailwind CSS for styling. It provides utility-first classes that help keep the HTML clean and maintainable.
