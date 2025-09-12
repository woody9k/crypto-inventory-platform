# GitHub Secrets Verification Checklist

## How to Test Your GitHub Secrets

### 1. Go to Your GitHub Repository
- Navigate to: `https://github.com/woody9k/crypto-inventory-platform`
- Click **Settings** (in the repository menu)
- Click **Secrets and variables** → **Actions**

### 2. Verify Required Secrets Exist

Check that these secrets are present:

#### ✅ Core Application Secrets (Required)
- [ ] `JWT_SECRET` - JWT signing secret
- [ ] `DATABASE_URL` - Test database connection
- [ ] `REDIS_URL` - Test Redis connection

#### ✅ Optional Integration Secrets
- [ ] `ATLASSIAN_API_TOKEN` - Atlassian API token
- [ ] `ATLASSIAN_SITE_URL` - Your Atlassian site URL
- [ ] `ATLASSIAN_EMAIL` - Your Atlassian email
- [ ] `NOTION_TOKEN` - Notion integration token
- [ ] `NOTION_DATABASE_ID` - Notion database ID
- [ ] `LINEAR_API_KEY` - Linear API key

### 3. Test the CI/CD Pipeline

#### Option A: Use the Existing CI Workflow
1. Go to **Actions** tab in your repository
2. Click on **CI/CD Pipeline**
3. Click **Run workflow** → **Run workflow**
4. Watch the workflow run and check for errors

#### Option B: Create a Simple Test (Manual)
1. Go to **Actions** tab
2. Click **New workflow**
3. Copy this test workflow:

```yaml
name: Test Secrets
on: [workflow_dispatch]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - name: Test JWT_SECRET
      run: |
        if [ -n "$JWT_SECRET" ]; then
          echo "✅ JWT_SECRET is set"
        else
          echo "❌ JWT_SECRET is missing"
          exit 1
        fi
    - name: Test DATABASE_URL
      run: |
        if [ -n "$DATABASE_URL" ]; then
          echo "✅ DATABASE_URL is set"
        else
          echo "❌ DATABASE_URL is missing"
          exit 1
        fi
    - name: Test REDIS_URL
      run: |
        if [ -n "$REDIS_URL" ]; then
          echo "✅ REDIS_URL is set"
        else
          echo "❌ REDIS_URL is missing"
          exit 1
        fi
    - name: Success
      run: echo "🎉 All secrets are working!"
```

### 4. Expected Results

If your secrets are configured correctly, you should see:
- ✅ All required secrets are set
- ✅ The workflow runs without errors
- ✅ All tests pass

### 5. Common Issues and Solutions

#### Issue: "Secret not found"
- **Solution**: Go to Settings → Secrets and add the missing secret

#### Issue: "Database connection failed"
- **Solution**: Check that `DATABASE_URL` is correct for the test environment

#### Issue: "Redis connection failed"
- **Solution**: Check that `REDIS_URL` is correct for the test environment

### 6. Next Steps After Verification

Once your secrets are working:
1. ✅ Your CI/CD pipeline will work correctly
2. ✅ Automated testing will run on every push
3. ✅ You can deploy with confidence
4. ✅ We can clean up the git history issue

## Quick Test Commands

You can also test locally that your secrets would work:

```bash
# Test database connection (if you have psql)
psql "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable" -c "SELECT 1;"

# Test Redis connection (if you have redis-cli)
redis-cli -h localhost -p 6379 ping
```

## Success Criteria

- [ ] All required secrets are present in GitHub
- [ ] CI/CD pipeline runs successfully
- [ ] No authentication errors in the workflow
- [ ] All services can connect to their dependencies
