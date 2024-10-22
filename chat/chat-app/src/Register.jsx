import React, { useState } from "react";
import { useNavigate } from "react-router-dom"; // 导入 useNavigate

const Register = () => {
  const navigate = useNavigate(); // 创建 navigate 函数
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");

  const handleRegisterSubmit = async (e) => {
    e.preventDefault();

    const registerData = {
      username: username,
      password: password,
    };

    try {
      const response = await fetch("http://localhost:8080/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(registerData),
      });

      const data = await response.json();
      if (response.ok) {
        setMessage("User registered successfully");
      } else {
        setMessage(data.error || "Registration failed");
      }
    } catch (error) {
      setMessage("An error occurred");
    }
  };

  return (
    <div>
      <h1>Register</h1>
      <form onSubmit={handleRegisterSubmit}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">Register</button>
      </form>
      {message && <p>{message}</p>}
      {/* 返回按钮 */}
      <button onClick={() => navigate("/")}>Back to Login</button>
    </div>
  );
};

export default Register;
