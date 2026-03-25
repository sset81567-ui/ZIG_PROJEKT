package main

import "fmt"

func GetLayout(title string, u *User) string {
 clr := "#007AFF"
 if u.ThemeColor != "" { clr = u.ThemeColor }

 return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ZIG | %s</title>
    <style>
        :root {
            --accent: %s;
            --bg: #000000;
            --panel: #0a0a0a;
            --text: #ffffff;
            --gray: #8e8e93;
            --border: #1c1c1e;
        }
        * { box-sizing: border-box; }
        body { 
            background: var(--bg); color: var(--text); 
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto; 
            margin: 0; display: flex; height: 100vh; overflow: hidden;
        }
        
        /* Sidebar */
        .side { 
            width: 380px; background: var(--panel); 
            border-right: 1px solid var(--border); 
            display: flex; flex-direction: column; 
        }
        .head { padding: 30px 20px; font-size: 28px; font-weight: 800; color: var(--accent); letter-spacing: -1.5px; }
        
        /* Profile Card */
        .card { 
            background: #151517; margin: 15px; padding: 25px; 
            border-radius: 28px; border: 1px solid var(--border);
            box-shadow: 0 10px 30px rgba(0,0,0,0.5);
        }
        .av-big { 
            width: 70px; height: 70px; background: var(--accent); 
            border-radius: 22px; display: flex; align-items: center; 
            justify-content: center; font-size: 30px; font-weight: bold; margin-bottom: 15px;
        }
        
        /* Stats/Gifts */
        .stats { display: flex; gap: 10px; margin: 15px 0; }
        .stat-item { background: #1c1c1e; padding: 8px 15px; border-radius: 12px; font-size: 13px; color: var(--gray); }
        .stat-item b { color: var(--text); }

        /* Inputs */
        input, textarea { 
            width: 100%%; background: #2c2c2e; border: 1px solid transparent; 
            color: #fff; padding: 14px; border-radius: 14px; margin-bottom: 12px;
            font-size: 15px; transition: 0.2s;
        }
        input:focus { border-color: var(--accent); outline: none; background: #3a3a3c; }

        .save-btn { 
            background: var(--accent); color: #fff; border: none; 
            width: 100%%; padding: 16px; border-radius: 16px; 
            font-weight: 700; cursor: pointer; transition: 0.3s;
        }
        .save-btn:hover { transform: translateY(-2px); box-shadow: 0 5px 15px var(--accent); }

        .main { flex-grow: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; }
        .logo-bg { font-size: 120px; font-weight: 900; opacity: 0.03; user-select: none; }
    </style>
</head>
<body>
    <div class="side">
        <div class="head">ZIG GLOBAL</div>
        <div class="card">
            <div class="av-big">Z</div>
            <h2 style="margin:0;">%s</h2>
            <div class="stats">
                <div class="stat-item">🧸 Мишки: <b>%d</b></div>
                <div class="stat-item">💎 Pro: <b>Нет</b></div>
            </div>
            <form action="/update" method="POST">
                <input type="text" name="username" placeholder="Никнейм" value="%s">
                <textarea name="bio" rows="3" placeholder="О себе...">%s</textarea>
                <button type="submit" class="save-btn">Сохранить профиль</button>
            </form>
        </div>
    </div>
    <div class="main">
        <div class="logo-bg">ZIG</div>
        <p style="color: var(--gray);">Выберите чат, чтобы начать общение</p>
    </div>
</body>
</html>`, title, clr, u.FullName, u.MishkaCount, u.Username, u.Bio)
}