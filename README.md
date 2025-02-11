# Bokkoli: A Veggie-to-Veggie Messaging CLI App 🌱

[![LinkedIn][linkedin-shield]][linkedin-url]

Bokkoli is a peer-to-peer **TCP messaging CLI application** built with **Go**, the **Bubbletea** framework, and **SQLite** for message persistence. The name "Bokkoli" is a playful fusion of "broccoli" and "보끔" (boggeum), the Korean word for "mix," reflecting its nature as a blend of technologies.

This lightweight chat application enables **real-time, direct peer-to-peer communication** through a clean, terminal-based interface. Users can easily connect and exchange messages using the TCP protocol, while **SQLite** ensures that conversation history is safely stored locally.

---

## Built With

- ![Go](./assets/Bokkoli-Golang-image.png)  
  **Go** for the backend and message handling.
  
- ![SQLite](./assets/Bokkoli-sqlite-image.png)  
  **SQLite** for persistent local message storage.
  
- ![Bubbletea](https://avatars.githubusercontent.com/u/24752999?s=200&v=4)  
  **Bubbletea** for the UI framework to create a clean terminal interface.

---

## Wireframe

Here’s a look at the design architecture:

**Internal Design:**  
<img src="./assets/Bokkoli-internalDesign.drawio.png" width="300" />

**Deployment Flow:**  
<img src="./assets/Bokkoli-Deployment.drawio.png" width="300" />

---



### Planned Features 🚀:

- [ ] **Fuzzy Matching**
- [ ] **Questions, Comments, and Suggestion Forum** within the app
- [ ] **Peer Authentication**
- [ ] **Handshake Protocol**
- [ ] **Message Integrity**
  - [ ] Checksums
  - [ ] Digital Signatures
- [ ] **Encryption and Security**
  - [ ] Data Encryption
  - [ ] End-to-End Encryption
- [ ] **Peer Status Monitoring**
  - [ ] Ping/Pong Mechanism
  - [ ] Health Checks
- [ ] **File/Content Sharing**
  - [ ] File Chunking
  - [ ] File Integrity Checks
- [ ] **Unit Testing** for reliability and coverage

---

## Installation 🔧

To get started with **Bokkoli**, follow
