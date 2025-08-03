sequenceDiagram
  participant Client
  participant API Server
  participant Database

  %% 1. Registration
  Client->>API Server: POST /register (email, password)
  API Server->>Database: Create user (hash password)
  Database-->>API Server: user_id
  API Server-->>Client: 201 Created

  %% 2. Login
  Client->>API Server: POST /login (email, password)
  API Server->>Database: Verify credentials
  alt valid credentials
    API Server->>Database: Create auth_session (hashed refresh token, device info)
    Database-->>API Server: session_id
    API Server-->>Client: access_token, refresh_token
  else invalid credentials
    API Server-->>Client: 401 Unauthorized
  end

  %% 3. Access Protected Resource
  Client->>API Server: GET /profile (Authorization: Bearer access_token)
  API Server->>API Server: Validate JWT (signature + expiry)
  alt valid JWT
    API Server-->>Client: 200 OK + Data
  else expired/invalid
    API Server-->>Client: 401 Unauthorized
  end

  %% 4. Refresh Token
  Client->>API Server: POST /refresh-token (refresh_token)
  API Server->>Database: Look up session (hash match)
  alt session found and valid
    API Server->>Database: Update last_used_at, rotate token
    API Server-->>Client: new access_token + refresh_token
  else not found / revoked / expired
    API Server-->>Client: 401 Unauthorized
  end

  %% 5. Logout (Invalidate Token)
  Client->>API Server: POST /logout (Authorization: Bearer refresh_token)
  API Server->>Database: Invalidate session (is_valid = false)
  API Server-->>Client: 204 No Content
