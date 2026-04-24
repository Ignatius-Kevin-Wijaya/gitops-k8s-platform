# AKS Terraform Module

## What this does

* Creates a Resource Group
* Provisions an AKS cluster
* Uses system-assigned managed identity
* Creates a default system node pool

## What this does NOT do (intentionally)

* No custom VNet
* No ACR integration
* No autoscaling
* No remote backend

---

## Usage

### 1. Create local variables file

Copy the example file to a local file named `terraform.tfvars`:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your own values.

---

### 2. Init

```bash
terraform init
```

---

### 3. Validate

```bash
terraform validate
```

---

### 4. Plan

```bash
terraform plan
```

---

### 5. Apply

```bash
terraform apply
```

---

## Notes

* Requires Azure login via Azure CLI
* Uses default networking (kubenet)
