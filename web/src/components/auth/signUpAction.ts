import { type ActionFunctionArgs, redirect } from "react-router";
import type { AppClient } from "../../lib/client";

export function signUpAction(client: AppClient) {
  return async ({ request }: ActionFunctionArgs) => {
    const formData = await request.formData();

    const username = formData.get("username")?.toString();
    const password = formData.get("password")?.toString();
    const confirmPassword = formData.get("confirmPassword")?.toString();

    const errors: Record<string, string> = {};

    if (!username) {
      errors.username = "Field is mandatory";
    } else if (username.length < 6) {
      errors.username = "Username should be at least 6 characters";
    }

    if (!password) {
      errors.password = "Field is mandatory";
    } else if (password.length < 6) {
      errors.password = "Password should be at least 6 characters";
    }

    if (!confirmPassword) {
      errors.confirmPassword = "Field is mandatory";
    }

    if (password && confirmPassword && password !== confirmPassword) {
      errors.confirmPassword = "Passwords do not match";
    }

    if (Object.keys(errors).length > 0) {
      return errors;
    }

    if (!username || !password) {
      return errors;
    }

    const safeUsername = username;
    const safePassword = password;

    const mutation = async () => {
      console.log("[signup] sending request", { username: safeUsername });
      return client.fetchClient.POST("/auth/register", {
        body: { username: safeUsername, password: safePassword },
      });
    };

    try {
      const res = await mutation();
      if (res.error) {
        console.error("[signup] api error");
        return { form: "Registration failed" };
      }
      console.log("[signup] registration succeeded", {
        username: safeUsername,
      });
    } catch (err) {
      console.error("[signup] network/client error", err);
      return { form: "Registration failed" };
    }

    return redirect("/signin");
  };
}
