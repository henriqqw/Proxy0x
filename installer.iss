#define MyAppName "Proxy0x"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "Proxy0x"
#define MyAppURL "https://github.com/henriqqw/Proxy0x"
#define MyAppExeName "Proxy0x.exe"

[Setup]
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
brazilianportuguese.LaunchAtStartup=Iniciar Proxy0x ao ligar o Windows
portuguese.LaunchAtStartup=Iniciar Proxy0x ao arrancar o Windows
spanish.LaunchAtStartup=Iniciar Proxy0x al arrancar Windows
french.LaunchAtStartup=Lancer Proxy0x au demarrage de Windows
german.LaunchAtStartup=Proxy0x beim Windows-Start starten
italian.LaunchAtStartup=Avvia Proxy0x all'avvio di Windows
russian.LaunchAtStartup=Zapuskat' Proxy0x pri zapuske Windows
japanese.LaunchAtStartup=Windows kidou ji ni Proxy0x wo kidou suru
korean.LaunchAtStartup=Windows sijaksi Proxy0x siljang
arabic.LaunchAtStartup=Tashghil Proxy0x ind bida' tashghil Windows
turkish.LaunchAtStartup=Windows basladiginda Proxy0x'i baslat
thai.LaunchAtStartup=Roem Proxy0x mua Windows roem thamngan
polish.LaunchAtStartup=Uruchom Proxy0x przy starcie Windows

english.CreateDesktopShortcut=Create a desktop shortcut
brazilianportuguese.CreateDesktopShortcut=Criar atalho na area de trabalho
portuguese.CreateDesktopShortcut=Criar atalho no ambiente de trabalho
spanish.CreateDesktopShortcut=Crear acceso directo en el escritorio
french.CreateDesktopShortcut=Creer un raccourci sur le bureau
german.CreateDesktopShortcut=Desktop-Verknuepfung erstellen
italian.CreateDesktopShortcut=Crea collegamento sul desktop
russian.CreateDesktopShortcut=Sozdat' yarlyk na rabochem stole
japanese.CreateDesktopShortcut=Desukutoppu ni shotokattu wo sakusei
korean.CreateDesktopShortcut=Bataseu hwamyeon-e barogagi mandeulgi
arabic.CreateDesktopShortcut=Insha' ikhtisaar ala sath al-maktab
turkish.CreateDesktopShortcut=Masaustune kisayol olustur
thai.CreateDesktopShortcut=Sangthang lad bon natha chob
polish.CreateDesktopShortcut=Utworz skrot na pulpicie

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopShortcut}"; GroupDescription: "{cm:AdditionalIcons}"
Name: "launchstartup"; Description: "{cm:LaunchAtStartup}"

[Files]
Source: "build\bin\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Registry]
Root: HKCU; Subkey: "SOFTWARE\Microsoft\Windows\CurrentVersion\Run"; ValueType: string; ValueName: "{#MyAppName}"; ValueData: """{app}\{#MyAppExeName}"""; Flags: uninsdeletevalue; Tasks: launchstartup

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#MyAppName}}"; Flags: nowait postinstall skipifsilent
