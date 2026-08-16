!!! info "Driver signing"
    Since [v.0.9.7.5](https://github.com/vadimgrn/usbip-win2/releases/tag/v.0.9.7.5) the driver is
    signed by Microsoft, so it loads with Secure Boot enabled and needs neither test signing nor an
    extra root CA. v.0.9.7.5 is WHQL-certified for x64 (passed Microsoft's HLK tests);
    [v.0.9.7.7](https://github.com/vadimgrn/usbip-win2/releases/tag/v.0.9.7.7) moved to
    [attestation signing](https://learn.microsoft.com/en-us/windows-hardware/drivers/dashboard/driver-signing-offerings)
    instead, a lighter Microsoft-trusted tier that skips HLK testing, and reports `Microsoft Windows
    Hardware Compatibility Publisher` as its signer. Both still load under Secure Boot on Windows
    10/11 Desktop. Signing is done through [OSSign](https://github.com/OSSign), so the installer
    itself carries their EV certificate.

    Older releases used to add the publicly available test signing CA as a _trusted root CA_. Checked
    on Windows 11 (build 26200) with Secure Boot enabled in August 2026: installing v.0.9.7.7
    (`USBip-0.9.7.7-x64.exe`, SHA256
    `51620FA5F9F8BE5932BC9D786DEEE557CE06D5407A99CAB490DCFAC71F185FEA`) added no certificates to
    `LocalMachine\Root`, `LocalMachine\TrustedPublisher` or `CurrentUser\Root`. If you're on an older
    release, compare the stores yourself around the install:

    ```powershell
    $when = 'before'  # switch to 'after' and re-run once the installer finishes
    $stores = 'Cert:\LocalMachine\Root','Cert:\LocalMachine\TrustedPublisher','Cert:\CurrentUser\Root'
    $stores | ForEach-Object { $s=$_; Get-ChildItem $s |
        Select-Object @{n='Store';e={$s}},Thumbprint,Subject } |
      Sort-Object Store,Thumbprint | Export-Csv "$env:USERPROFILE\certs-$when.csv" -NoTypeInformation
    ```

    Then `Compare-Object` the two CSVs. A test signing CA, if you find one, can be removed with
    `certmgr.msc` (as admin) by deleting "USBIP" from "Trusted Root Certification Authorities" ->
    "Certificates"; the driver still loads without it.

!!! warning "Avoid v.0.9.7.8"
    Its author [reports memory corruption and BSODs](https://github.com/vadimgrn/usbip-win2/releases/tag/v.0.9.7.8)
    in that release. As of August 2026 nothing newer has shipped, so
    [v.0.9.7.7](https://github.com/vadimgrn/usbip-win2/releases/tag/v.0.9.7.7) is the newest release
    without such a report; it also fixes a BSOD from accessing paged memory at high IRQL that was
    present in earlier releases.
