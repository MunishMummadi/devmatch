"use client"

import { useEffect, useState } from "react"

export default function MessagesPage() {
    const [conversations, setConversations] = useState([])
    const [messages, setMessages] = useState([])
    const [selectedConversation, setSelectedConversation] = useState(null)

    useEffect(() => {
        const fetchData = async () => {
            const res = await fetch("/mock/chatData.json")
            const data = await res.json()

            // Add unread count to conversations
            const conversationsWithUnread = data.conversations.map((conv, index) => ({
                ...conv,
                unread: index % 3 === 0 ? Math.floor(Math.random() * 5) + 1 : 0, // Add random unread count to some conversations
            }))

            setConversations(conversationsWithUnread)

            if (conversationsWithUnread.length > 0) {
                const firstId = conversationsWithUnread[0].id
                setSelectedConversation(firstId)
                setMessages(data.messages[firstId] || [])
            }
        }
        fetchData()
    }, [])

    const handleConversationClick = (id) => {
        setSelectedConversation(id)
        fetch("/mock/chatData.json")
            .then((res) => res.json())
            .then((data) => {
                setMessages(data.messages[id] || [])
            })
    }

    return (
        <div style={styles.container}>
            {/* Left column */}
            <div style={styles.conversationList}>
                <div style={styles.appLogo}>
                    <span style={styles.logoIcon}>❤️</span> SOUL
                </div>
                <div style={styles.searchBar}>
                    <input type="text" placeholder="Search matches" style={styles.searchInput} />
                </div>

                {conversations.map((conv) => (
                    <div
                        key={conv.id}
                        onClick={() => handleConversationClick(conv.id)}
                        style={{
                            ...styles.conversationItem,
                            backgroundColor: selectedConversation === conv.id ? "#fff8fa" : "#fff",
                            borderLeft: selectedConversation === conv.id ? "4px solid #a01a40" : "none",
                        }}
                    >
                        <img src={conv.avatar || "/placeholder.svg"} alt={conv.contactName} style={styles.avatar} />
                        <div style={styles.conversationDetails}>
                            <div style={styles.conversationTitle}>
                                {conv.contactName}
                                <span style={styles.timeStamp}>{conv.timestamp}</span>
                            </div>
                            <div style={styles.conversationDesc}>{conv.title}</div>
                            <p style={styles.preview}>{conv.lastMessage}</p>
                        </div>
                        {conv.unread && <div style={styles.unreadBadge}>{conv.unread}</div>}
                    </div>
                ))}
            </div>

            {/* Right column */}
            <div style={styles.chatArea}>
                {selectedConversation && (
                    <>
                        <div style={styles.chatHeader}>
                            {conversations.find((c) => c.id === selectedConversation) && (
                                <>
                                    <img
                                        src={conversations.find((c) => c.id === selectedConversation).avatar || "/placeholder.svg"}
                                        alt="Profile"
                                        style={styles.chatHeaderAvatar}
                                    />
                                    <div>
                                        <div style={styles.chatHeaderName}>
                                            {conversations.find((c) => c.id === selectedConversation).contactName}
                                        </div>
                                        <div style={styles.chatHeaderStatus}>Online</div>
                                    </div>
                                </>
                            )}
                        </div>
                        <div style={styles.chatMessages}>
                            {messages.map((msg, index) => (
                                <div
                                    key={index}
                                    style={{
                                        ...styles.messageItem,
                                        ...(msg.type === "outgoing" && styles.messageItemMine),
                                    }}
                                >
                                    <div style={msg.type === "outgoing" ? styles.bubbleMine : styles.bubble}>
                                        <p style={{ margin: "0 0 0.5rem 0" }}>{msg.text}</p>
                                        <small
                                            style={{
                                                ...styles.chatTime,
                                                ...(msg.type !== "outgoing" && styles.chatTimeTheirs),
                                            }}
                                        >
                                            {msg.timestamp}
                                        </small>
                                    </div>
                                </div>
                            ))}
                        </div>

                        <div style={styles.chatFooter}>
                            <div style={styles.attachButton}>📎</div>
                            <input type="text" placeholder="Type a message..." style={styles.chatInput} />
                            <div style={styles.emojiButton}>😊</div>
                            <div style={styles.sendButton}>➤</div>
                        </div>
                    </>
                )}
            </div>
        </div>
    )
}

const styles = {
    container: {
        display: "flex",
        height: "100vh",
        backgroundColor: "#f8f0f2",
    },
    conversationList: {
        width: "320px",
        backgroundColor: "#fff",
        display: "flex",
        flexDirection: "column",
        borderRight: "1px solid #e0e0e0",
        boxShadow: "0 0 15px rgba(0,0,0,0.05)",
    },
    searchBar: {
        padding: "1rem",
        borderBottom: "1px solid #f0e0e5",
        backgroundColor: "#a01a40",
    },
    searchInput: {
        width: "100%",
        padding: "0.7rem 1.2rem",
        borderRadius: "24px",
        border: "none",
        outline: "none",
        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
        fontSize: "0.9rem",
    },
    conversationItem: {
        display: "flex",
        gap: "1rem",
        padding: "1rem",
        borderBottom: "1px solid #f0e0e5",
        cursor: "pointer",
        transition: "all 0.2s ease",
    },
    avatar: {
        width: "50px",
        height: "50px",
        borderRadius: "50%",
        objectFit: "cover",
        border: "2px solid #f0f0f0",
        boxShadow: "0 2px 5px rgba(0,0,0,0.1)",
    },
    conversationDetails: {
        flex: 1,
    },
    conversationTitle: {
        fontWeight: 600,
        display: "flex",
        justifyContent: "space-between",
        fontSize: "1rem",
        marginBottom: "0.3rem",
        color: "#333",
    },
    timeStamp: {
        fontWeight: 400,
        fontSize: "0.8rem",
        color: "#999",
    },
    conversationDesc: {
        fontSize: "0.9rem",
        color: "#555",
        marginBottom: "0.3rem",
    },
    preview: {
        fontSize: "0.85rem",
        color: "#888",
        marginTop: 0,
        whiteSpace: "nowrap",
        overflow: "hidden",
        textOverflow: "ellipsis",
        maxWidth: "220px",
    },
    chatArea: {
        flex: 1,
        display: "flex",
        flexDirection: "column",
        backgroundColor: "#fff",
        position: "relative",
    },
    chatHeader: {
        display: "flex",
        alignItems: "center",
        padding: "1rem",
        borderBottom: "1px solid #f0e0e5",
        backgroundColor: "#a01a40",
        color: "white",
    },
    chatHeaderAvatar: {
        width: "40px",
        height: "40px",
        borderRadius: "50%",
        marginRight: "1rem",
        border: "2px solid white",
    },
    chatHeaderName: {
        fontWeight: 600,
        fontSize: "1.1rem",
    },
    chatHeaderStatus: {
        fontSize: "0.8rem",
        opacity: 0.8,
    },
    chatMessages: {
        flex: 1,
        padding: "1.5rem",
        overflowY: "auto",
        backgroundImage: "linear-gradient(to bottom, #fff8fa, #fff)",
    },
    messageItem: {
        display: "flex",
        alignItems: "flex-start",
        marginBottom: "1.2rem",
    },
    messageItemMine: {
        justifyContent: "flex-end",
    },
    bubble: {
        backgroundColor: "#f5f5f5",
        borderRadius: "18px 18px 18px 0",
        padding: "0.8rem 1.2rem",
        maxWidth: "70%",
        boxShadow: "0 1px 2px rgba(0,0,0,0.1)",
    },
    bubbleMine: {
        backgroundColor: "#a01a40",
        color: "#fff",
        borderRadius: "18px 18px 0 18px",
        padding: "0.8rem 1.2rem",
        maxWidth: "70%",
        boxShadow: "0 1px 2px rgba(0,0,0,0.1)",
    },
    chatTime: {
        display: "block",
        fontSize: "0.75rem",
        color: "rgba(255,255,255,0.7)",
        marginTop: "0.4rem",
        textAlign: "right",
    },
    chatTimeTheirs: {
        color: "#999",
    },
    chatFooter: {
        display: "flex",
        alignItems: "center",
        padding: "1rem 1.2rem",
        borderTop: "1px solid #f0e0e5",
        backgroundColor: "white",
    },
    chatInput: {
        flex: 1,
        padding: "0.8rem 1.2rem",
        borderRadius: "24px",
        border: "1px solid #e0e0e0",
        outline: "none",
        marginRight: "1rem",
        fontSize: "0.95rem",
        boxShadow: "0 1px 3px rgba(0,0,0,0.05)",
    },
    chatFooterIcons: {
        display: "flex",
        gap: "1rem",
        fontSize: "1.2rem",
        color: "#a01a40",
    },
    sendButton: {
        backgroundColor: "#a01a40",
        color: "white",
        width: "40px",
        height: "40px",
        borderRadius: "50%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        cursor: "pointer",
        boxShadow: "0 2px 5px rgba(160,26,64,0.3)",
    },
    attachButton: {
        color: "#a01a40",
        cursor: "pointer",
    },
    emojiButton: {
        color: "#a01a40",
        cursor: "pointer",
    },
    unreadBadge: {
        backgroundColor: "#a01a40",
        color: "white",
        borderRadius: "50%",
        width: "20px",
        height: "20px",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        fontSize: "0.7rem",
        marginLeft: "auto",
    },
    appLogo: {
        display: "flex",
        alignItems: "center",
        padding: "1rem",
        backgroundColor: "#a01a40",
        color: "white",
        fontWeight: "bold",
        fontSize: "1.2rem",
    },
    logoIcon: {
        marginRight: "0.5rem",
    },
}
