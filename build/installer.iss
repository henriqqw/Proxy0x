; Saved as UTF-8 with BOM so the non-Latin messages compile correctly.
; Build: iscc /DMyAppVersion=1.2.3 build\installer.iss  (run from the repo root or anywhere)
#define MyAppName "Proxy0x"
#ifndef MyAppVersion
  #define MyAppVersion "0.0.0-dev"
#endif
#define MyAppPublisher "Proxy0x"
#define MyAppURL "https://github.com/henriqqw/Proxy0x"
#define MyAppExeName "Proxy0x.exe"
#define MyMCPExeName "proxy0x-mcp.exe"

[Setup]
; Paths below are relative to the repository root, not to this script's folder.
SourceDir=..
AppId={{B7E3F2A1-4D5C-4E6F-8A9B-0C1D2E3F4A5B}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
OutputDir=build\installer
OutputBaseFilename=Proxy0x-Setup-{#MyAppVersion}
SetupIconFile=build\windows\icon.ico
UninstallDisplayIcon={app}\{#MyAppExeName}
Compression=lzma2/ultra64
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
ArchitecturesInstallIn64BitMode=x64compatible
CloseApplications=yes
RestartApplications=no

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "brazilianportuguese"; MessagesFile: "compiler:Languages\BrazilianPortuguese.isl"
Name: "portuguese"; MessagesFile: "compiler:Languages\Portuguese.isl"
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"
Name: "french"; MessagesFile: "compiler:Languages\French.isl"
Name: "german"; MessagesFile: "compiler:Languages\German.isl"
Name: "italian"; MessagesFile: "compiler:Languages\Italian.isl"
Name: "russian"; MessagesFile: "compiler:Languages\Russian.isl"
Name: "japanese"; MessagesFile: "compiler:Languages\Japanese.isl"
Name: "korean"; MessagesFile: "compiler:Languages\Korean.isl"
Name: "arabic"; MessagesFile: "compiler:Languages\Arabic.isl"
Name: "turkish"; MessagesFile: "compiler:Languages\Turkish.isl"
Name: "thai"; MessagesFile: "compiler:Languages\Thai.isl"
Name: "polish"; MessagesFile: "compiler:Languages\Polish.isl"

[CustomMessages]
english.LaunchAtStartup=Start Proxy0x when Windows starts
brazilianportuguese.LaunchAtStartup=Iniciar o Proxy0x junto com o Windows
portuguese.LaunchAtStartup=Iniciar o Proxy0x ao arrancar o Windows
spanish.LaunchAtStartup=Iniciar Proxy0x al arrancar Windows
french.LaunchAtStartup=Lancer Proxy0x au démarrage de Windows
german.LaunchAtStartup=Proxy0x beim Windows-Start starten
italian.LaunchAtStartup=Avvia Proxy0x all'avvio di Windows
russian.LaunchAtStartup=Запускать Proxy0x при запуске Windows
japanese.LaunchAtStartup=Windows の起動時に Proxy0x を起動する
korean.LaunchAtStartup=Windows 시작 시 Proxy0x 실행
arabic.LaunchAtStartup=تشغيل Proxy0x عند بدء تشغيل Windows
turkish.LaunchAtStartup=Windows başladığında Proxy0x'i başlat
thai.LaunchAtStartup=เปิด Proxy0x เมื่อเริ่ม Windows
polish.LaunchAtStartup=Uruchamiaj Proxy0x przy starcie systemu Windows

english.CreateDesktopShortcut=Create a desktop shortcut
brazilianportuguese.CreateDesktopShortcut=Criar atalho na área de trabalho
portuguese.CreateDesktopShortcut=Criar atalho no ambiente de trabalho
spanish.CreateDesktopShortcut=Crear acceso directo en el escritorio
french.CreateDesktopShortcut=Créer un raccourci sur le bureau
german.CreateDesktopShortcut=Desktop-Verknüpfung erstellen
italian.CreateDesktopShortcut=Crea collegamento sul desktop
russian.CreateDesktopShortcut=Создать ярлык на рабочем столе
japanese.CreateDesktopShortcut=デスクトップにショートカットを作成する
korean.CreateDesktopShortcut=바탕 화면에 바로 가기 만들기
arabic.CreateDesktopShortcut=إنشاء اختصار على سطح المكتب
turkish.CreateDesktopShortcut=Masaüstüne kısayol oluştur
thai.CreateDesktopShortcut=สร้างทางลัดบนเดสก์ท็อป
polish.CreateDesktopShortcut=Utwórz skrót na pulpicie

english.RemoveUserData=Also delete saved proxies and settings (AppData\Roaming\Proxy0x)?
brazilianportuguese.RemoveUserData=Apagar também os proxies e as configurações salvos (AppData\Roaming\Proxy0x)?
portuguese.RemoveUserData=Apagar também os proxies e as definições guardados (AppData\Roaming\Proxy0x)?
spanish.RemoveUserData=¿Eliminar también los proxies y la configuración guardados (AppData\Roaming\Proxy0x)?
french.RemoveUserData=Supprimer aussi les proxys et les paramètres enregistrés (AppData\Roaming\Proxy0x) ?
german.RemoveUserData=Gespeicherte Proxys und Einstellungen ebenfalls löschen (AppData\Roaming\Proxy0x)?
italian.RemoveUserData=Eliminare anche i proxy e le impostazioni salvati (AppData\Roaming\Proxy0x)?
russian.RemoveUserData=Удалить также сохранённые прокси и настройки (AppData\Roaming\Proxy0x)?
japanese.RemoveUserData=保存されたプロキシと設定 (AppData\Roaming\Proxy0x) も削除しますか?
korean.RemoveUserData=저장된 프록시와 설정(AppData\Roaming\Proxy0x)도 삭제하시겠습니까?
arabic.RemoveUserData=هل تريد أيضًا حذف الوكلاء والإعدادات المحفوظة (AppData\Roaming\Proxy0x)؟
turkish.RemoveUserData=Kayıtlı proxy'ler ve ayarlar da silinsin mi (AppData\Roaming\Proxy0x)?
thai.RemoveUserData=ลบพร็อกซีและการตั้งค่าที่บันทึกไว้ (AppData\Roaming\Proxy0x) ด้วยหรือไม่?
polish.RemoveUserData=Usunąć także zapisane proxy i ustawienia (AppData\Roaming\Proxy0x)?

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopShortcut}"; GroupDescription: "{cm:AdditionalIcons}"
Name: "launchstartup"; Description: "{cm:LaunchAtStartup}"

[Files]
Source: "build\bin\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "build\bin\{#MyMCPExeName}"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Registry]
Root: HKCU; Subkey: "SOFTWARE\Microsoft\Windows\CurrentVersion\Run"; ValueType: string; ValueName: "{#MyAppName}"; ValueData: """{app}\{#MyAppExeName}"""; Flags: uninsdeletevalue; Tasks: launchstartup

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#MyAppName}}"; Flags: nowait postinstall skipifsilent

[Code]
// netsh exits with 0 when the rule exists and 1 when it does not.
function KillSwitchActive(): Boolean;
var
  ResultCode: Integer;
begin
  Result := Exec(ExpandConstant('{sys}\netsh.exe'),
    'advfirewall firewall show rule name=Proxy0x-KillSwitch-AllowApp',
    '', SW_HIDE, ewWaitUntilTerminated, ResultCode) and (ResultCode = 0);
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
var
  ResultCode: Integer;
begin
  // Runs as the user so the policy backup in their AppData is found; the app
  // asks for administrator rights itself.
  if (CurUninstallStep = usUninstall) and KillSwitchActive() then
    Exec(ExpandConstant('{app}\{#MyAppExeName}'), '--disable-killswitch',
      '', SW_HIDE, ewWaitUntilTerminated, ResultCode);

  if (CurUninstallStep = usPostUninstall) and not UninstallSilent then
    if MsgBox(CustomMessage('RemoveUserData'), mbConfirmation, MB_YESNO or MB_DEFBUTTON2) = IDYES then
      DelTree(ExpandConstant('{userappdata}\Proxy0x'), True, True, True);
end;
