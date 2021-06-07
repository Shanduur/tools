$trigger = New-JobTrigger -AtStartup -RandomDelay 00:00:30
Register-ScheduledJob -Trigger $trigger -FilePath D:\Scripts\tools\jiggler\run.ps1 -Name Run-Jiggler