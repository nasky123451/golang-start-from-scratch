import React, { useState, useEffect, useRef } from 'react';
import { jwtDecode } from 'jwt-decode';
import {
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemText,
  Typography,
  TextField,
  Button,
  Box,
  Divider,
  Card,
} from '@mui/material';
import { Menu as MenuIcon, Close as CloseIcon } from '@mui/icons-material';

const Chat = () => {
  const [currentUser, setCurrentUser] = useState('');
  const [onlineUsers, setOnlineUsers] = useState([]);
  const [offlineUsers, setOfflineUsers] = useState([]);
  const [messages, setMessages] = useState([]);
  const [messageInput, setMessageInput] = useState('');
  const [ws, setWs] = useState(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const messageEndRef = useRef(null);
  const chatContainerRef = useRef(null);
  const [noMoreMessages, setNoMoreMessages] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      try {
        const decoded = jwtDecode(token);
        setCurrentUser(decoded.username);
        fetchOnlineUsers();
        fetchMessages(new Date().toISOString().split('T')[0], "general", false);
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

  useEffect(() => {
    console.log(messages);
  }, [messages]);

  const fetchOnlineUsers = async () => {
    try {
      const response = await fetch('http://localhost:8080/online-users', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await response.json();
      setOnlineUsers(data.onlineUsers || []);
      setOfflineUsers(data.offlineUsers || []);
    } catch (error) {
      console.error('Failed to fetch online users:', error);
    }
  };

  const fetchMessages = async (date, room, separator) => {
    try {
      const response = await fetch(`http://localhost:8080/chat-history?date=${date || new Date().toISOString().split('T')[0]}&room=${room}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();

      if (Array.isArray(data.messages)) {
        if (data.messages.length === 0) {
          console.warn('No messages found for the selected date.');
          setNoMoreMessages(true);
        } else {
          setNoMoreMessages(false);

          if (separator) {
            // Create date separator message
            const previousDate = new Date(); // Get current date
            previousDate.setDate(previousDate.getDate() - 1); // Get the previous day's date

            // Add date separator message to message list
            const dateSeparatorMessage = {
                isSeparator: true,
                date: previousDate.toISOString(),
            };

            // Merge new messages with previous messages
            setMessages((prevMessages) => [
                ...data.messages,
                dateSeparatorMessage,
                ...prevMessages
            ]);
          } else {
            setMessages((prevMessages) => [
                ...data.messages,
                ...prevMessages
            ]);
          }

          // Use setTimeout to ensure that the scroll position is checked after the message is loaded
          setTimeout(() => {
            chatContainerRef.current.scrollTop += messageEndRef.current.clientHeight; // Adjust scroll position
          }, 0);
        }
      } else {
        console.error('Messages is not an array or is undefined:', data.messages);
      }
    } catch (error) {
      console.error('Failed to fetch messages:', error);
    }
  };

  const setupWebSocket = (username) => {
    const token = localStorage.getItem('token');
    const socket = new WebSocket(`ws://localhost:8080/ws`);

    socket.onopen = () => {
      socket.send(JSON.stringify({ type: "auth", token }));
      setWs(socket);
    };

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);

      if (msg.type === "message") {
        setMessages((prevMessages) => [...prevMessages, msg]);
        // Automatically scroll only when at the bottom
        if (chatContainerRef.current.scrollHeight - chatContainerRef.current.scrollTop === chatContainerRef.current.clientHeight) {
          scrollToBottom();
        }
      } else if (msg.type === "userStatus") {
        updateUserStatus(msg.username, msg.status);
      }
    };

    socket.onclose = () => {
      console.log('WebSocket connection closed');
    };
  };

  const updateUserStatus = (username, status) => {
    if (status === 'online') {
      setOnlineUsers((prev) => [...new Set([...prev, username])]);
      setOfflineUsers((prev) => prev.filter(user => user !== username));
    } else if (status === 'offline') {
      setOfflineUsers((prev) => [...new Set([...prev, username])]);
      setOnlineUsers((prev) => prev.filter(user => user !== username));
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

    ws.send(JSON.stringify(message));
    setMessageInput('');
  };

  const scrollToBottom = () => {
    if (messageEndRef.current) {
        messageEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
};

  const handleScroll = (e) => {
    const { scrollTop } = e.target;

    if (scrollTop === 0) {
      // Save the current number of message lists
      const currentMessagesCount = messages.length;

      const lastMessageDate = messages.length > 0 ? new Date(messages[1].time) : null;
      const previousDate = lastMessageDate ? new Date(lastMessageDate) : new Date();
      previousDate.setDate(previousDate.getDate() - 1);
      fetchMessages(previousDate.toISOString().split('T')[0], "general", true);

      // Use setTimeout to ensure that the scroll position is checked after the message is loaded
      setTimeout(() => {
        // If the number of new messages loaded does not increase, keep the original scroll position
        if (messages.length === currentMessagesCount) {
          chatContainerRef.current.scrollTop = scrollTop; // Restore previous scroll position
        }
      }, 100);
    }
  };

  // Make sure to scroll to the bottom when loading
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    scrollToBottom();
  }, []);

  return (
    <div style={{ display: 'flex', height: '100vh' }}>
      {/* Drawer components */}
      <Drawer
        anchor="left"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        variant="temporary"
      >
        <Box sx={{ width: 250, padding: '20px', backgroundColor: '#f0f0f0' }}>
          <IconButton onClick={() => setDrawerOpen(false)} sx={{ marginBottom: '16px' }}>
            <CloseIcon />
          </IconButton>
          <Typography variant="h6" gutterBottom>當前用戶</Typography>
          <Typography variant="body1" sx={{ marginBottom: '16px' }}>{currentUser}</Typography>
          <Divider sx={{ marginBottom: '16px' }} />
          <Typography variant="h6" gutterBottom>在線用戶</Typography>
          <List>
            {(onlineUsers || [])
              .filter(user => user !== currentUser)
              .map((user, index) => (
                <ListItem key={index} sx={{ '&:hover': { backgroundColor: '#e0e0e0' } }}>
                  <ListItemText primary={user} />
                </ListItem>
              ))}
          </List>
          <Divider sx={{ marginBottom: '16px' }} />
          <Typography variant="h6" gutterBottom>離線用戶</Typography>
          <List>
            {(offlineUsers || [])
              .filter(user => user !== currentUser)
              .map((user, index) => (
                <ListItem key={index} sx={{ '&:hover': { backgroundColor: '#e0e0e0' } }}>
                  <ListItemText primary={user} />
                </ListItem>
              ))}
          </List>
        </Box>
      </Drawer>

      {/* main content */}
      <div style={{ flexGrow: 1, padding: '20px' }}>
        <IconButton onClick={() => setDrawerOpen(true)}>
          <MenuIcon />
        </IconButton>

        <Typography variant="h4" gutterBottom align="center">聊天窗口</Typography>

        <Card sx={{ mb: 2, padding: 2 }}>
          <Box 
            ref={chatContainerRef}
            sx={{ 
              height: '400px', 
              overflowY: 'scroll', 
              border: '1px solid #ccc', 
              padding: '10px', 
              backgroundColor: '#f9f9f9' 
            }}
            onScroll={handleScroll}
          >
            {noMoreMessages && (
              <Typography variant="body2" color="textSecondary" align="center" style={{ marginBottom: '10px' }}>
                没有更多消息了
              </Typography>
            )}
            {(messages || []).map((msg, index) => (
              <div key={index}>
                {msg.isSeparator ? (
                  <Typography variant="body2" align="center" sx={{ margin: '10px 0', fontWeight: 'bold' }}>
                    {new Date(msg.date).toLocaleDateString()} {/* Show date */}
                  </Typography>
                ) : (
                  <Box sx={{ marginBottom: 1 }}>
                    <strong>{msg.from}:</strong> {msg.content} <em>{new Date(msg.time).toLocaleTimeString()}</em>
                  </Box>
                )}
              </div>
            ))}
            <div ref={messageEndRef} />
          </Box>
        </Card>

        <form onSubmit={sendMessage}>
          <TextField 
            value={messageInput}
            onChange={(e) => setMessageInput(e.target.value)}
            label="輸入消息..."
            fullWidth
            variant="outlined"
            sx={{ mb: 1 }}
          />
          <Button type="submit" variant="contained" color="primary">發送</Button>
        </form>
      </div>
    </div>
  );
};

export default Chat;
