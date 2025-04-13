// app/messages/page.js

export default function MessagesPage() {
    return (
        <div style={styles.container}>
            {/* LEFT SIDEBAR with icons */}
            <aside style={styles.sidebar}>
                <div style={styles.profileArea}>
                    {/* Placeholder user avatar */}
                    <img
                        src="https://via.placeholder.com/40"
                        alt="User Avatar"
                        style={styles.avatar}
                    />
                </div>
                <div style={styles.navIcons}>
                    <div style={styles.navIcon}>
                        <span>🏠</span>
                    </div>
                    <div style={styles.navIcon}>
                        <span>✉️</span>
                    </div>
                    <div style={styles.navIcon}>
                        <span>📂</span>
                    </div>
                    <div style={styles.navIcon}>
                        <span>🗑️</span>
                    </div>
                </div>
                <button style={styles.fabButton}>+</button>
            </aside>

            {/* CONVERSATION LIST */}
            <div style={styles.conversationList}>
                <div style={styles.searchBar}>
                    <input
                        type="text"
                        placeholder="Search"
                        style={styles.searchInput}
                    />
                </div>
                {/* Hard-coded conversation items */}
                <div style={styles.conversationItem}>
                    <img
                        src="https://via.placeholder.com/40"
                        alt="Jessica Carroll"
                        style={styles.listAvatar}
                    />
                    <div style={styles.conversationDetails}>
                        <div style={styles.conversationTitle}>
                            Jessica Carroll
                            <span style={styles.timeStamp}>1h ago</span>
                        </div>
                        <div style={styles.conversationDesc}>
                            Work Inquiry - UI Designer
                        </div>
                        <p style={styles.conversationPreview}>
                            Currently we are looking for a UI designer...
                        </p>
                    </div>
                </div>

                <div style={styles.conversationItem}>
                    <img
                        src="https://via.placeholder.com/40"
                        alt="Emily Rose"
                        style={styles.listAvatar}
                    />
                    <div style={styles.conversationDetails}>
                        <div style={styles.conversationTitle}>
                            Emily Rose
                            <span style={styles.timeStamp}>2h ago</span>
                        </div>
                        <div style={styles.conversationDesc}>
                            Invitation: Board Game Night
                        </div>
                        <p style={styles.conversationPreview}>
                            Sed rhoncus aliquam velit sit amet...
                        </p>
                    </div>
                </div>

                <div style={styles.conversationItem}>
                    <img
                        src="https://via.placeholder.com/40"
                        alt="David Bryant"
                        style={styles.listAvatar}
                    />
                    <div style={styles.conversationDetails}>
                        <div style={styles.conversationTitle}>
                            David Bryant
                            <span style={styles.timeStamp}>5h ago</span>
                        </div>
                        <div style={styles.conversationDesc}>
                            New App Design
                        </div>
                        <p style={styles.conversationPreview}>
                            Lorem ipsum dolor sit amet...
                        </p>
                    </div>
                </div>
            </div>

            {/* CHAT WINDOW */}
            <main style={styles.chatArea}>
                <header style={styles.chatHeader}>
                    <div>
                        <h3 style={styles.chatTitle}>Jessica Carroll</h3>
                        <small style={styles.onlineStatus}>Online</small>
                    </div>
                    <div style={styles.chatActions}>
                        <span>📞</span>
                        <span>📱</span>
                        <span>⋮</span>
                    </div>
                </header>

                <div style={styles.chatMessages}>
                    <div style={styles.messageItem}>
                        <img
                            src="https://via.placeholder.com/40"
                            alt="Jessica Carroll"
                            style={styles.listAvatar}
                        />
                        <div style={styles.bubble}>
                            <p>Hey Michael!</p>
                            <small style={styles.chatTime}>10:17 am</small>
                        </div>
                    </div>

                    <div style={{ ...styles.messageItem, ...styles.messageItemMine }}>
                        <div style={styles.bubbleMine}>
                            <p>Hello Jessica, how can I help you?</p>
                            <small style={styles.chatTime}>10:29 am</small>
                        </div>
                    </div>

                    <div style={styles.messageItem}>
                        <img
                            src="https://via.placeholder.com/40"
                            alt="Jessica Carroll"
                            style={styles.listAvatar}
                        />
                        <div style={styles.bubble}>
                            <p>
                                Currently we are looking for a UI designer to work on our websites
                                and mobile application.
                            </p>
                            <small style={styles.chatTime}>10:35 am</small>
                        </div>
                    </div>
                </div>

                <footer style={styles.chatFooter}>
                    <input
                        type="text"
                        placeholder="Type your message here"
                        style={styles.chatInput}
                    />
                    <div style={styles.chatFooterIcons}>
                        <span>📎</span>
                        <span>😊</span>
                        <span>➤</span>
                    </div>
                </footer>
            </main>
        </div>
    );
}

const styles = {
    container: {
        display: 'flex',
        height: '100vh',
        backgroundColor: '#f5f7fa',
    },
    // LEFT SIDEBAR
    sidebar: {
        width: '60px',
        backgroundColor: '#2196f3',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        padding: '1rem 0',
        gap: '1rem',
    },
    profileArea: {
        marginBottom: '1rem',
    },
    avatar: {
        width: '40px',
        height: '40px',
        borderRadius: '50%',
        objectFit: 'cover',
    },
    navIcons: {
        display: 'flex',
        flexDirection: 'column',
        gap: '1rem',
        flex: 1,
    },
    navIcon: {
        cursor: 'pointer',
        color: '#fff',
        textAlign: 'center',
        fontSize: '1.2rem',
    },
    fabButton: {
        backgroundColor: '#ff9800',
        border: 'none',
        borderRadius: '50%',
        width: '40px',
        height: '40px',
        fontSize: '1.2rem',
        color: '#fff',
        cursor: 'pointer',
    },

    // CONVERSATION LIST
    conversationList: {
        width: '300px',
        backgroundColor: '#fff',
        display: 'flex',
        flexDirection: 'column',
        borderRight: '1px solid #e0e0e0',
    },
    searchBar: {
        padding: '1rem',
        borderBottom: '1px solid #e0e0e0',
    },
    searchInput: {
        width: '100%',
        padding: '0.5rem 1rem',
        borderRadius: '20px',
        border: '1px solid #ccc',
        outline: 'none',
    },
    conversationItem: {
        display: 'flex',
        gap: '0.8rem',
        padding: '1rem',
        borderBottom: '1px solid #e0e0e0',
        cursor: 'pointer',
    },
    listAvatar: {
        width: '40px',
        height: '40px',
        borderRadius: '50%',
    },
    conversationDetails: {
        flex: 1,
    },
    conversationTitle: {
        fontWeight: 600,
        display: 'flex',
        justifyContent: 'space-between',
        fontSize: '0.95rem',
        marginBottom: '0.2rem',
    },
    timeStamp: {
        fontWeight: 400,
        fontSize: '0.8rem',
        color: '#777',
    },
    conversationDesc: {
        fontSize: '0.85rem',
        color: '#555',
        marginBottom: '0.3rem',
    },
    conversationPreview: {
        fontSize: '0.8rem',
        color: '#999',
        marginTop: 0,
    },

    // CHAT AREA
    chatArea: {
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
    },
    chatHeader: {
        display: 'flex',
        justifyContent: 'space-between',
        padding: '1rem',
        borderBottom: '1px solid #e0e0e0',
        backgroundColor: '#fff',
    },
    chatTitle: {
        margin: 0,
        fontSize: '1.1rem',
    },
    onlineStatus: {
        color: 'green',
        fontSize: '0.85rem',
        marginLeft: '0.5rem',
    },
    chatActions: {
        display: 'flex',
        gap: '1rem',
        fontSize: '1.2rem',
        alignItems: 'center',
    },
    chatMessages: {
        flex: 1,
        padding: '1rem',
        overflowY: 'auto',
    },
    messageItem: {
        display: 'flex',
        alignItems: 'flex-start',
        marginBottom: '1rem',
        gap: '0.8rem',
    },
    messageItemMine: {
        flexDirection: 'row-reverse',
    },
    bubble: {
        backgroundColor: '#eee',
        borderRadius: '10px',
        padding: '0.8rem 1rem',
        maxWidth: '60%',
    },
    bubbleMine: {
        backgroundColor: '#2196f3',
        color: '#fff',
        borderRadius: '10px',
        padding: '0.8rem 1rem',
        maxWidth: '60%',
    },
    chatTime: {
        display: 'block',
        fontSize: '0.75rem',
        color: '#777',
        marginTop: '0.4rem',
    },
    chatFooter: {
        display: 'flex',
        alignItems: 'center',
        padding: '0.8rem 1rem',
        borderTop: '1px solid #e0e0e0',
        backgroundColor: '#fff',
    },
    chatInput: {
        flex: 1,
        padding: '0.5rem 1rem',
        borderRadius: '20px',
        border: '1px solid #ccc',
        outline: 'none',
        marginRight: '1rem',
    },
    chatFooterIcons: {
        display: 'flex',
        gap: '1rem',
        fontSize: '1.2rem',
        cursor: 'pointer',
    },
};
