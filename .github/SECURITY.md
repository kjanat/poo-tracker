# Security Policy

## 🔒 Supported Versions

We actively support the following versions of Poo Tracker with security updates:

| Version | Supported |
| ------- | --------- |
| 1.x.x   | ✅ Yes    |
| < 1.0   | ❌ No     |

## 🚨 Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please follow these steps:

### 🔐 Private Disclosure

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please report security issues through one of these methods:

1. **GitHub Security Advisories** (Preferred)

   - Go to [Security Advisories](https://github.com/kjanat/poo-tracker/security/advisories/new)
   - Click "Report a vulnerability"
   - Fill out the form with details

2. **Email**

   - Send details to: <poo-tracker+security@kjanat.com>
   - Include "SECURITY" in the subject line
   - Use PGP encryption if possible (key available on request)

3. **Discord** (For urgent issues)
   - Direct message @kjanat on Discord
   - Mention it's a security issue

### 📋 What to Include

When reporting a security vulnerability, please include:

- **Description**: Clear description of the vulnerability
- **Impact**: What could an attacker accomplish?
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Proof of Concept**: Code snippets, screenshots, or videos if applicable
- **Affected Components**: Which parts of the system are affected?
- **Suggested Fix**: If you have ideas for how to fix it

### ⏱️ Response Timeline

- **Initial Response**: Within 24 hours
- **Assessment**: Within 72 hours
- **Status Updates**: Every 7 days until resolved
- **Fix Release**: Depends on severity (see below)

### 🚦 Severity Levels

| Severity     | Response Time | Description                                     | Examples                 |
| ------------ | ------------- | ----------------------------------------------- | ------------------------ |
| **Critical** | 24-48 hours   | Remote code execution, authentication bypass    | SQL injection, RCE       |
| **High**     | 3-7 days      | Significant data exposure, privilege escalation | Data leaks, admin access |
| **Medium**   | 1-2 weeks     | Limited impact, requires specific conditions    | XSS, CSRF                |
| **Low**      | 2-4 weeks     | Minimal impact, theoretical vulnerabilities     | Information disclosure   |

## 🏆 Recognition

We believe in recognizing security researchers who help make Poo Tracker safer:

- **Hall of Fame**: Contributors listed in our security acknowledgments
- **Swag**: Poo Tracker stickers and merchandise for significant findings
- **References**: We'll provide LinkedIn recommendations for professional researchers

## 🛡️ Security Measures

### Current Security Practices

- **Authentication**: JWT tokens with proper expiration
- **Authorization**: Role-based access control
- **Data Protection**: Encryption at rest and in transit
- **Input Validation**: Comprehensive input sanitization
- **Dependency Management**: Regular security audits via Dependabot
- **Container Security**: Regular base image updates
- **CI/CD Security**: Automated security scanning in pipelines

### Infrastructure Security

- **Database**: PostgreSQL with encrypted connections
- **File Storage**: S3-compatible storage with proper IAM
- **API Security**: Rate limiting, CORS policies, security headers
- **Docker**: Non-root containers, minimal attack surface
- **Secrets**: Environment-based secret management

## 🚫 Out of Scope

The following are generally considered out of scope:

- **Self-hosted instances**: Security of user-managed deployments
- **Third-party integrations**: Issues in external services we integrate with
- **Social engineering**: Attacks requiring human interaction
- **Physical security**: Physical access to servers
- **DDoS attacks**: Availability attacks (report to infrastructure team)
- **Brute force attacks**: On properly rate-limited endpoints

## 📚 Security Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Node.js Security Best Practices](https://nodejs.org/en/docs/guides/security/)
- [React Security Guidelines](https://react.dev/learn/security)
- [Docker Security Best Practices](https://docs.docker.com/engine/security/)

## 🤝 Responsible Disclosure

We follow responsible disclosure practices:

1. **Private reporting** of vulnerabilities
2. **Coordination** on timeline for fixes
3. **Public disclosure** only after fixes are released
4. **Credit** given to researchers (unless requested otherwise)

## 💡 Security Tips for Users

- Keep your deployment updated to the latest version
- Use strong authentication credentials
- Enable HTTPS in production
- Regularly backup your data
- Monitor your logs for suspicious activity
- Use environment variables for sensitive configuration

## 📞 Contact

For security-related questions or concerns:

- **Security Team**: <poo-tracker+security@kjanat.com>
- **General Questions**: Use GitHub Discussions
- **Urgent Issues**: Discord @kjanat

---

Remember: Security is everyone's responsibility. When in doubt, report it! 🛡️
