# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to: **security@terraship.io**

You should receive a response within 48 hours. If for some reason you do not, please follow up via email to ensure we received your original message.

Please include the following information:

- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

This information will help us triage your report more quickly.

## Security Best Practices

When using Terraship:

### Credentials Management

- **Never commit credentials** to version control
- Use environment variables or credential files outside the repository
- Use IAM roles when running in cloud environments
- Rotate credentials regularly
- Follow the principle of least privilege

### Policy Files

- Review policy files before use
- Keep policies in version control
- Use strict validation rules for production
- Regularly update policies as infrastructure evolves

### CI/CD Integration

- Use encrypted secrets for cloud credentials
- Limit GitHub Action permissions using `permissions:` key
- Use OIDC/federated credentials when possible
- Review action logs for sensitive information exposure

### Network Security

- Run Terraship in private networks when possible
- Use VPC endpoints for cloud API access
- Enable audit logging for Terraship executions
- Monitor for unusual validation patterns

## Known Security Considerations

### Terraform State Files

Terraship reads Terraform state which may contain sensitive information. Ensure:

- State files are encrypted at rest
- Access to state is properly controlled
- State backends use encryption in transit
- State file access is audited

### Cloud Provider Access

Terraship requires read access to cloud resources. To minimize risk:

- Use read-only IAM policies where possible
- Scope permissions to specific resources
- Use separate credentials for validation
- Monitor API usage for anomalies

### Policy Evaluation

- Policies are executed in the Terraship process
- Complex regex patterns may cause performance issues
- Validate policy files before deployment
- Test policies in non-production environments first

## Disclosure Policy

When we receive a security bug report, we will:

1. Confirm the problem and determine affected versions
2. Audit code to find similar problems
3. Prepare fixes for all supported versions
4. Release patches as soon as possible

We will credit researchers who report vulnerabilities responsibly.

## Security Updates

Security updates will be released as:

- Patch versions for minor issues
- Minor versions for significant issues
- Emergency releases for critical vulnerabilities

Subscribe to our [security advisories](https://github.com/vijayaxai/terraship/security/advisories) to receive notifications.

## Comments on this Policy

If you have suggestions on how this process could be improved, please submit a pull request.

---

Last updated: January 2026
