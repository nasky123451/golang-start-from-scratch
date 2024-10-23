import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Container,
  TextField,
  Button,
  Typography,
  Box,
  Snackbar,
  Alert,
} from "@mui/material";

const Register = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [openSnackbar, setOpenSnackbar] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false); // 新增狀態來追蹤註冊是否成功

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
        setIsSuccess(true); // 設置成功狀態
        setOpenSnackbar(true); // 打开 Snackbar
      } else {
        setMessage(data.error || "Registration failed");
        setIsSuccess(false); // 設置失敗狀態
        setOpenSnackbar(true); // 打开 Snackbar
      }
    } catch (error) {
      console.error("Error:", error); // 打印错误以帮助调试
      setMessage("An error occurred");
      setIsSuccess(false); // 設置失敗狀態
      setOpenSnackbar(true); // 打开 Snackbar
    }
  };

  const handleCloseSnackbar = () => {
    setOpenSnackbar(false);
  };

  return (
    <Container maxWidth="xs" sx={{ mt: 5 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Register
      </Typography>
      <form onSubmit={handleRegisterSubmit}>
        <Box display="flex" flexDirection="column" gap={2}>
          <TextField
            label="Username"
            variant="outlined"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <TextField
            label="Password"
            type="password"
            variant="outlined"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <Button type="submit" variant="contained" color="primary">
            Register
          </Button>
          <Button
            onClick={() => navigate("/")}
            variant="outlined"
            color="secondary"
          >
            Back to Login
          </Button>
        </Box>
      </form>
      <Snackbar open={openSnackbar} autoHideDuration={6000} onClose={handleCloseSnackbar}>
        <Alert onClose={handleCloseSnackbar} severity={isSuccess ? "success" : "error"}>
          {message}
        </Alert>
      </Snackbar>
    </Container>
  );
};

export default Register;
