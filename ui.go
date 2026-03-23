package main

import "fmt"

func ui(color string, isChat bool) string {
 if color == "" { color = "#0088cc" }
 sD, cD := "flex", "none"
 if isChat { sD, cD = "none", "flex" }

 return fmt.Sprintf(`
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
    :root { --main: %s; --bg: #0e1621; --side: #17212b; --text: #fff; }
    body { background: var(--bg); color: var(--text); font-family: sans-serif; margin: 0; display: flex; height: 100vh; overflow: hidden; }
    .sidebar { width: 320px; background: var(--side); border-right: 1px solid #000; display: %s; flex-direction: column; }
    .main-chat { flex: 1; display: %s; flex-direction: column; background: #070d14; }
    @media (min-width: 769px) { .sidebar, .main-chat { display: flex !important; } }
    .chat-item { padding: 15px; display: flex; align-items: center; text-decoration: none; color: #fff; border-bottom: 1px solid #0e1621; }
    .avatar { width: 40px; height: 40px; background: var(--main); border-radius: 50%%; margin-right: 12px; display: flex; align-items: center; justify-content: center; }
    .messages { flex: 1; padding: 20px; overflow-y: auto; display: flex; flex-direction: column; gap: 10px; }
    .bubble { max-width: 75%%; padding: 10px; border-radius: 12px; background: #182533; }
    .bubble.me { align-self: flex-end; background: var(--main); }
    .input-bar { padding: 10px; background: var(--side); display: flex; gap: 10px; }
    input { flex: 1; padding: 12px; border-radius: 20px; border: none; background: #242f3d; color: #fff; outline: none; }
    button { background: var(--main); color: #fff; border: none; padding: 10px 20px; border-radius: 20px; cursor: pointer; }
</style>
<script>
    setInterval(() => {
        const chat = new URLSearchParams(window.location.search).get('chat');
        if (chat) {
            fetch('/api/messages?chat=' + chat)
                .then(r => r.text())
                .then(html => {
                    const box = document.getElementById('msg-box');
                    if (box && box.innerHTML.trim() !== html.trim()) {
                        box.innerHTML = html;
                        box.scrollTop = box.scrollHeight;
                    }
                });
        }
    }, 2000);
</script>`, color, sD, cD)
}