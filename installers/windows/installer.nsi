; ═══════════════════════════════════════════════════════════
; خان (Khan) — Windows Installer (NSIS)
; چت سازمانی سبک فارسی
; ═══════════════════════════════════════════════════════════

Unicode true
SetCompressor /SOLID lzma

; ─── اطلاعات برنامه ───
!define APP_NAME "خان"
!define APP_NAME_EN "Khan Chat"
!define APP_VERSION "1.0.3"
!define APP_PUBLISHER "Khan Team"
!define APP_EXE "khan.exe"
!define APP_URL "https://codeberg.org/adiib/khan1.0.2"

; ─── مسیرها ───
!define APP_DIR "$PROGRAMFILES64\Khan"
!define APP_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\Khan"
!define UNINST_EXE "$INSTDIR\Uninstall.exe"

Name "${APP_NAME}"
OutFile "KhanSetup-${APP_VERSION}.exe"
InstallDir "${APP_DIR}"
InstallDirRegKey HKLM "${APP_UNINST_KEY}" "InstallLocation"
RequestExecutionLevel admin
BrandingText "Khan Chat v${APP_VERSION}"

; ─── صفحات ───
Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

; ─── زبان (فارسی/عربی RTL + انگلیسی) ───
!include "MUI2.nsh"
!insertmacro MUI_LANGUAGE "Arabic"
!insertmacro MUI_LANGUAGE "English"

; ─── بخش نصب ───
Section "Install"
  SectionIn RO
  SetOutPath "${APP_DIR}"
  
  ; فایل اصلی
  File "${APP_EXE}"
  
  ; پوشه داده
  CreateDirectory "${APP_DIR}\data"
  
  ; میان‌بر دسکتاپ
  CreateShortcut "$DESKTOP\${APP_NAME}.lnk" "${APP_DIR}\${APP_EXE}" "" "${APP_DIR}\${APP_EXE}"
  
  ; میان‌بر منوی استارت
  CreateDirectory "$SMPROGRAMS\${APP_NAME}"
  CreateShortcut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" "${APP_DIR}\${APP_EXE}"
  CreateShortcut "$SMPROGRAMS\${APP_NAME}\حذف ${APP_NAME}.lnk" "$INSTDIR\Uninstall.exe"
  
  ; رجیستری
  WriteRegStr HKLM "${APP_UNINST_KEY}" "DisplayName" "${APP_NAME} — چت سازمانی"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "DisplayVersion" "${APP_VERSION}"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "Publisher" "${APP_PUBLISHER}"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "DisplayIcon" "${APP_DIR}\${APP_EXE}"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "UninstallString" "$INSTDIR\Uninstall.exe"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "InstallLocation" "${APP_DIR}"
  WriteRegDWORD HKLM "${APP_UNINST_KEY}" "NoModify" 1
  WriteRegDWORD HKLM "${APP_UNINST_KEY}" "NoRepair" 1
  
  ; یونیستالر
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  
  ; ─── فایروال ویندوز: باز کردن پورت 1727 (برای دسترسی LAN) ───
  ; حذف rule قدیمی (اگر هست) و ساخت rule جدید
  nsExec::ExecToLog 'netsh advfirewall firewall delete rule name="Khan Chat 1727"'
  nsExec::ExecToLog 'netsh advfirewall firewall add rule name="Khan Chat 1727" dir=in action=allow protocol=TCP localport=1727 profile=private'
  nsExec::ExecToLog 'netsh advfirewall firewall add rule name="Khan Chat 1727" dir=in action=allow protocol=TCP localport=1727 profile=domain'
SectionEnd

; ─── بخش حذف ───
Section "Uninstall"
  ; میان‌برها
  Delete "$DESKTOP\${APP_NAME}.lnk"
  Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
  Delete "$SMPROGRAMS\${APP_NAME}\حذف ${APP_NAME}.lnk"
  RMDir "$SMPROGRAMS\${APP_NAME}"
  
  ; فایل‌ها (داده حفظ می‌شود — بکاپ!)
  Delete "${APP_DIR}\\${APP_EXE}"
  
  ; حذف rule فایروال
  nsExec::ExecToLog 'netsh advfirewall firewall delete rule name="Khan Chat 1727"'
  
  ; رجیستری
  DeleteRegKey HKLM "${APP_UNINST_KEY}"
  
  ; پوشه
  RMDir "${APP_DIR}"
SectionEnd
