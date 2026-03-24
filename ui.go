package main

import "fmt"

func GetLayout(content string, u *User) string {
 // Настройки по умолчанию
 theme := "#3498db"
 lang := "ru"
 bio := "Привет! Я использую ZIG."
 name := "Новый пользователь"
 
 if u != nil {
  if u.ThemeColor != "" { theme = u.ThemeColor }
  if u.Language != "" { lang = u.Language }
  if u.Bio != "" { bio = u.Bio }
  name = u.GetDisplayName()
 }

 return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="%s">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>ZIG Global</title>
    <style>
        :root {
            --main-clr: %s;
            --bg-black: #000000;
            --panel-bg: #0a0a0a;
            --card-bg: #111111;
            --border-clr: #1f1f1f;
            --text-gray: #888888;
        }

        * { box-sizing: border-box; -webkit-tap-highlight-color: transparent; }
        
        body { 
            background: var(--bg-black); 
            color: #ffffff; 
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; 
            margin: 0; 
            display: flex; 
            height: 100vh; 
            overflow: hidden; 
        }

        /* Боковая панель (Sidebar) */
        .sidebar { 
            width: 380px; 
            background: var(--panel-bg); 
            border-right: 1px solid var(--border-clr); 
            display: flex; 
            flex-direction: column; 
            transition: all 0.3s ease;
        }

        .header { 
            padding: 24px 20px; 
            font-size: 28px; 
            font-weight: 900; 
            color: var(--main-clr); 
            letter-spacing: -1.5px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        /* Поиск */
        .search-area { padding: 0 15px 15px 15px; }
        .search-input { 
            width: 100%%; 
            padding: 14px 18px; 
            border-radius: 14px; 
            border: 1px solid var(--border-clr); 
            background: #141414; 
            color: #fff; 
            font-size: 16px; 
            outline: none;
        }
        .search-input:focus { border-color: var(--main-clr); background: #1a1a1a; }

        /* Секция профиля */
        .scroll-content { flex-grow: 1; overflow-y: auto; padding: 10px; }
        
        .profile-card { 
            background: var(--card-bg); 
            padding: 24px; 
            border-radius: 24px; 
            border: 1px solid var(--border-clr); 
            margin-bottom: 15px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.5);
        }

        .avatar-container {
            display: flex;
            align-items: center;
            gap: 15px;
            margin-bottom: 20px;
        }

        .big-avatar { 
            width: 64px; 
            height: 64px; 
            background: var(--main-clr); 
            border-radius: 22px; 
            display: flex; 
            align-items: center; 
            justify-content: center; 
            font-size: 28px; 
            font-weight: bold;
            box-shadow: 0 0 20px var(--main-clr);
        }

        .user-info b { font-size: 18px; display: block; }
        .user-info span { color: var(--text-gray); font-size: 13px; }

        /* Форма настроек */
        .field-label { 
            font-size: 11px; 
            color: var(--text-gray); 
            text-transform: uppercase; 
            font-weight: 800; 
            margin: 15px 0 6px 4px; 
            display: block; 
        }

        .input-box { 
            width: 100%%; 
            background: #1a1a1a; 
            border: 1px solid var(--border-clr); 
            color: #fff; 
            padding: 12px; 
            border-radius: 12px; 
            font-size: 15px;
            margin-bottom: 5px;
        }

        .color-well {
            width: 100%%;
            height: 45px;
            border: none;
            background: none;
            cursor: pointer;
            border-radius: 10px;
        }

        .save-btn { 
            background: var(--main-clr); 
            color: #fff; 
            border: none; 
            padding: 16px; 
            width: 100%%; 
            border-radius: 16px; 
            font-weight: bold; 
            font-size: 16px;
            cursor: pointer; 
            margin-top: 20px; 
            transition: 0.2s;
        }
        .save-btn:active { transform: scale(0.97); }

        /* Главная область */
        .chat-area { 
            flex-grow: 1; 
            display: flex; 
            flex-direction: column; 
            align-items: center; 
            justify-content: center; 
            background: #050505; 
            text-align: center;
        }

        .placeholder-icon { font-size: 80px; color: #1a1a1a; margin-bottom: 20px; }

        @media (max-width: 768px) {
            .sidebar { width: 100%%; }
            .chat-area { display: none; }
        }
    </style>
</head>
<body>
    <div class="sidebar">
        <div class="header">
            <span>ZIG</span>
            <span style="font-size: 12px; opacity: 0.3;">v1.0-Global</span>
        </div>
        
        <div class="search-area">
            <input type="text" class="search-input" placeholder="Поиск (4-24 символа)...">
        </div>

        <div class="scroll-content">
            <div class="profile-card">
                <div class="avatar-container">
                    <div class="big-avatar">Z</div>
                    <div class="user-info">
                        <b>%s</b>
                        <span>Zoom In Global</span>
                    </div>
                </div>

                <form action="/update-profile" method="POST">
                    <span class="field-label">О себе (Bio)</span>
                    <textarea name="bio" class="input-box" rows="3" placeholder="Расскажите о себе...">%s</textarea>
                    
                    <span class="field-label">Язык</span>
                    <select name="lang" class="input-box">
                        <option value="ru" %s>Русский (Russian)</option>
                        <option value="en" %s>English (UK)</option>
                    </select>

                    <span class="field-label">Цвет оформления</span>
                    <input type="color" name="theme" class="color-well" value="%s">
                    
                    <button type="submit" class="save-btn">Сохранить профиль</button>
                </form>
            </div>
        </div>
    </div>

    <div class="chat-area">
        <div class="placeholder-icon">ZIG</div>
        <h2 style="color: #222;">Выберите чат, чтобы начать общение</h2>
    </div>
</body>
</html>`, lang, theme, name, bio, 
 func()string{if lang=="ru"{return "selected"};return ""}(), 
 func()string{if lang=="en"{return "selected"};return ""}(), 
 theme)
}