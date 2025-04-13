# DevMatch - Developer Connection Platform
![alt text](image.png)

A full-stack platform for connecting developers based on GitHub activity and technical interests.

## Features
- GitHub profile synchronization
- Real-time messaging system
- Conversation tracking with rich message history
- Developer matching algorithm

## Tech Stack
**Backend**
- Go 1.21
- Gin Web Framework
- GORM (PostgreSQL)
- GitHub API integration

**Frontend**
- React 18
- TypeScript
- Axios for API calls
- Tailwind CSS

## Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Docker (optional)

### Installation
```bash
# Clone repository
git clone https://github.com/MunishMummadi/devmatch.git
cd devmatch

# Install backend dependencies
cd server
go mod download

# Install frontend dependencies
cd ../client
npm install
```

### Configuration
Create `.env` files:

**Backend (.env)**
```env
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_secret
DATABASE_URL=postgres://user:pass@localhost:5432/devmatch
API_PORT=8080
```

**Frontend (.env)**
```env
REACT_APP_API_BASE_URL=http://localhost:8080/api
```

### Running
```bash
# Start backend
cd server
go run main.go

# Start frontend
cd ../client
npm start
```

## API Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users/{id}` | Get user profile |
| POST | `/api/conversations` | Start new conversation |
| GET | `/api/conversations/{id}/messages` | Get message history |

## Collaborators 
- [Naga Pavan Muni](https://github.com/nagapavanmuni) -> Fronend Developer
- [Munish Mummadi](https://github.com/MunishMummadi) -> Backend Developer
- [Sowmya Koneti](https://github.com/SowmyaKoneti) -> Fullstack Developer

## GitHub Integration
The system automatically:
- Syncs user profiles every 24 hours
- Handles API rate limits with exponential backoff
- Stores GitHub metadata including public repos and follower data

## License
- [MIT](LICENSE)
