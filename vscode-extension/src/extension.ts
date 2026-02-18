import * as vscode from 'vscode';
import * as child_process from 'child_process';
import * as path from 'path';

export function activate(context: vscode.ExtensionContext) {
    console.log('Terraship extension is now active');

    // Register validate command
    let validateDisposable = vscode.commands.registerCommand('terraship.validate', async () => {
        await runValidation();
    });

    // Register validate file command
    let validateFileDisposable = vscode.commands.registerCommand('terraship.validateFile', async () => {
        const editor = vscode.window.activeTextEditor;
        if (editor) {
            const fileDir = path.dirname(editor.document.uri.fsPath);
            await runValidation(fileDir);
        }
    });

    // Register init command
    let initDisposable = vscode.commands.registerCommand('terraship.init', async () => {
        await runInit();
    });

    // Register on-save validation if enabled
    let onSaveDisposable = vscode.workspace.onDidSaveTextDocument(async (document) => {
        const config = vscode.workspace.getConfiguration('terraship');
        const validateOnSave = config.get<boolean>('validateOnSave');
        
        if (validateOnSave && document.languageId === 'terraform') {
            const fileDir = path.dirname(document.uri.fsPath);
            await runValidation(fileDir, true);
        }
    });

    context.subscriptions.push(validateDisposable, validateFileDisposable, initDisposable, onSaveDisposable);
}

async function runValidation(directory?: string, silent: boolean = false) {
    const config = vscode.workspace.getConfiguration('terraship');
    const policyPath = config.get<string>('policyPath') || './policies/sample-policy.yml';
    const mode = config.get<string>('mode') || 'validate-existing';
    const cloudProvider = config.get<string>('cloudProvider') || '';
    let executablePath = config.get<string>('executablePath') || 'terraship';

    const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
    if (!workspaceFolder) {
        vscode.window.showErrorMessage('No workspace folder open');
        return;
    }

    const targetDir = directory || workspaceFolder.uri.fsPath;

    // Improve executable path handling for Windows
    if (process.platform === 'win32' && !executablePath.endsWith('.exe')) {
        executablePath = executablePath + '.exe';
    }

    if (!silent) {
        vscode.window.showInformationMessage('Running Terraship validation...');
    }

    const args = [
        'validate',
        targetDir,
        '--policy', policyPath,
        '--mode', mode,
        '--output', 'json',
        '--output-file', path.join(targetDir, 'terraship-report.json')
    ];

    if (cloudProvider) {
        args.push('--provider', cloudProvider);
    }

    try {
        const result = await execAsync(executablePath, args, { cwd: workspaceFolder.uri.fsPath });
        
        // Read and parse the report
        const reportPath = path.join(targetDir, 'terraship-report.json');
        const fs = require('fs');
        
        if (fs.existsSync(reportPath)) {
            const reportData = fs.readFileSync(reportPath, 'utf8');
            const report = JSON.parse(reportData);
            
            if (report.failed_resources === 0 && report.error_resources === 0) {
                vscode.window.showInformationMessage(
                    `✓ Validation passed: ${report.passed_resources} resources validated`
                );
            } else {
                vscode.window.showWarningMessage(
                    `✗ Validation failed: ${report.failed_resources} failures, ${report.warning_resources} warnings`
                );
                
                // Show detailed results
                showValidationResults(report);
            }
        }
    } catch (error: any) {
        const errorMsg = error.message || String(error);
        
        // Provide helpful error message for ENOENT (executable not found)
        if (errorMsg.includes('ENOENT') || errorMsg.includes('not found')) {
            const configuredPath = config.get<string>('executablePath') || 'terraship';
            vscode.window.showErrorMessage(
                `Terraship executable not found: "${executablePath}"\n\n` +
                `Please configure the path in VS Code settings:\n` +
                `1. Open Settings (Ctrl+,)\n` +
                `2. Search for "terraship.executablePath"\n` +
                `3. Set it to the full path to terraship.exe or ensure it's in your PATH`,
                'Open Settings'
            ).then(selection => {
                if (selection === 'Open Settings') {
                    vscode.commands.executeCommand('workbench.action.openSettings', 'terraship.executablePath');
                }
            });
        } else {
            vscode.window.showErrorMessage(`Terraship validation failed: ${errorMsg}`);
        }
    }
}

async function runInit() {
    const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
    if (!workspaceFolder) {
        vscode.window.showErrorMessage('No workspace folder open');
        return;
    }

    const config = vscode.workspace.getConfiguration('terraship');
    let executablePath = config.get<string>('executablePath') || 'terraship';

    // Improve executable path handling for Windows
    if (process.platform === 'win32' && !executablePath.endsWith('.exe')) {
        executablePath = executablePath + '.exe';
    }

    vscode.window.showInformationMessage('Initializing Terraship policy...');

    try {
        await execAsync(executablePath, ['init'], { cwd: workspaceFolder.uri.fsPath });
        vscode.window.showInformationMessage('Terraship policy initialized successfully');
    } catch (error: any) {
        const errorMsg = error.message || String(error);
        
        // Provide helpful error message for ENOENT (executable not found)
        if (errorMsg.includes('ENOENT') || errorMsg.includes('not found')) {
            const configuredPath = config.get<string>('executablePath') || 'terraship';
            vscode.window.showErrorMessage(
                `Terraship executable not found: "${executablePath}"\n\n` +
                `Please configure the path in VS Code settings:\n` +
                `1. Open Settings (Ctrl+,)\n` +
                `2. Search for "terraship.executablePath"\n` +
                `3. Set it to the full path to terraship.exe or ensure it's in your PATH`,
                'Open Settings'
            ).then(selection => {
                if (selection === 'Open Settings') {
                    vscode.commands.executeCommand('workbench.action.openSettings', 'terraship.executablePath');
                }
            });
        } else {
            vscode.window.showErrorMessage(`Failed to initialize policy: ${errorMsg}`);
        }
    }
}

function showValidationResults(report: any) {
    const panel = vscode.window.createWebviewPanel(
        'terrashipResults',
        'Terraship Validation Results',
        vscode.ViewColumn.Beside,
        {}
    );

    let html = '<html><body style="padding: 20px; font-family: sans-serif;">';
    html += '<h1>Terraship Validation Results</h1>';
    html += `<p><strong>Total Resources:</strong> ${report.total_resources}</p>`;
    html += `<p><strong>Passed:</strong> ${report.passed_resources}</p>`;
    html += `<p><strong>Failed:</strong> ${report.failed_resources}</p>`;
    html += `<p><strong>Warnings:</strong> ${report.warning_resources}</p>`;
    html += '<h2>Failed Resources</h2>';
    
    for (const resource of report.reports || []) {
        if (resource.status === 'fail') {
            html += `<div style="margin: 10px 0; padding: 10px; border-left: 3px solid red;">`;
            html += `<strong>${resource.resource_address}</strong> (${resource.resource_type})<br/>`;
            
            for (const result of resource.rule_results || []) {
                if (!result.passed) {
                    html += `<p style="margin: 5px 0;">✗ ${result.rule_name}: ${result.message}</p>`;
                    if (result.remediation) {
                        html += `<p style="margin: 5px 0; color: #666;">💡 ${result.remediation}</p>`;
                    }
                }
            }
            html += '</div>';
        }
    }
    
    html += '</body></html>';
    panel.webview.html = html;
}

function execAsync(command: string, args: string[], options: any): Promise<string> {
    return new Promise((resolve, reject) => {
        child_process.execFile(command, args, options, (error, stdout, stderr) => {
            if (error) {
                reject(error);
            } else {
                resolve(stdout.toString());
            }
        });
    });
}

export function deactivate() {}
