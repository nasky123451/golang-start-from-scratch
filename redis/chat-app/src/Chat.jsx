import React, { useState, useEffect } from 'react';
import { jwtDecode } from 'jwt-decode';

const Chat = () => {
  const [currentUser, setCurrentUser] = useState('');
  const [onlineUsers, setOnlineUsers] = useState([]); // 当前在线用户列表
  const [offlineUsers, setOfflineUsers] = useState([]); // 当前离线用户列表
  const [messages, setMessages] = useState([]);
  const [messageInput, setMessageInput] = useState('');
  const [ws, setWs] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      try {
        const decoded = jwtDecode(token);
        setCurrentUser(decoded.username);
        fetchOnlineUsers();
        setupWebSocket(decoded.username);
      } catch (error) {
        console.error("Token decoding error:", error);
      }
    }

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, []);

  const fetchOnlineUsers = async () => {
    try {
      const response = await fetch('http://localhost:8080/online-users', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await response.json();
      console.log('Fetched online users:', data); // Log the fetched data
      setOnlineUsers(data.onlineUsers || []); // Set to empty array if undefined
      setOfflineUsers(data.offlineUsers || []); // Set to empty array if undefined
    } catch (error) {
      console.error('Failed to fetch online users:', error);
    }
  };

  const setupWebSocket = (username) => {
    const token = localStorage.getItem('token');
    const socket = new WebSocket(`ws://localhost:8080/ws`);

    socket.onopen = () => {
      console.log('WebSocket connection established');
      socket.send(JSON.stringify({ type: "auth", token }));
      setWs(socket);
    };

    socket.onmessage = (event) => {
        console.log('Received message:', event.data);
      const msg = JSON.parse(event.data);
      console.log('Received message:', msg);

      if (msg.type === "message") {
        setMessages((prevMessages) => [...prevMessages, msg]);
      } else if (msg.type === "userStatus") {
        // 根据用户状态更新在线/离线用户
        updateUserStatus(msg.username, msg.status);
      }
    };

    socket.onclose = () => {
      console.log('WebSocket connection closed');
    };

  };

  const updateUserStatus = (username, status) => {
    if (status === 'online') {
      setOnlineUsers((prev) => [...new Set([...prev, username])]); // 添加在线用户
      setOfflineUsers((prev) => prev.filter(user => user !== username)); // 从离线用户中移除
    } else if (status === 'offline') {
      setOfflineUsers((prev) => [...new Set([...prev, username])]); // 添加离线用户
      setOnlineUsers((prev) => prev.filter(user => user !== username)); // 从在线用户中移除
    }
  };

  const sendMessage = (e) => {
    e.preventDefault();
    if (!messageInput || !ws) return;

    const message = {
      type: "message",
      room: 'general',
      from: currentUser,
      content: messageInput,
      time: new Date().toISOString(),
    };

    console.log('Sending message:', message);
    ws.send(JSON.stringify(message));
    setMessageInput('');
  };

  return (
    <div style={{ display: 'flex', height: '100vh' }}>
      <div style={{ width: '20%', borderRight: '1px solid #ccc' }}>
        <h3>当前用户</h3>
        <p>{currentUser}</p>
        <h4>在线用户</h4>
        <ul>
          {(onlineUsers || []).map((user, index) => (
            <li key={index}>{user}</li>
          ))}
        </ul>
        <h4>离线用户</h4>
        <ul>
          {(offlineUsers || []).map((user, index) => (
            <li key={index}>{user}</li>
          ))}
        </ul>
      </div>
      <div style={{ width: '80%', padding: '20px' }}>
        <h2>聊天窗口</h2>
        <div style={{ height: '400px', overflowY: 'scroll', border: '1px solid #ccc', padding: '10px' }}>
          {(messages || []).map((msg, index) => (
            <div key={index}>
              <strong>{msg.from}:</strong> {msg.content} <em>{new Date(msg.time).toLocaleTimeString()}</em>
            </div>
          ))}
        </div>
        <form onSubmit={sendMessage}>
          <input
            type="text"
            value={messageInput}
            onChange={(e) => setMessageInput(e.target.value)}
            placeholder="输入消息..."
            required
          />
          <button type="submit">发送</button>
        </form>
      </div>
    </div>
  );
};

export default Chat;
