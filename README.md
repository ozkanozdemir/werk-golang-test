# werk-golang-test

.env.example dosyasını kopyalayarak .env dosyası oluşturun ve içine hasura bilgilerini ekleyin

Hasura ve Go Server'ı için postman endpoint dokümanı dosya olarak eklenmiştir. Postman ile açıp hasura admin key'i header'a veya environment olarak ekleyerek istekleri atabilirsiniz.

### Graphql Sorguları

```
1- Şirket ekleme 

mutation {
  insert_companies(objects: [{
    name: "Company 1"
  }]) {
    returning {
      id
      name
      created_at
      updated_at
    }
  }
}


2- Kullanıcı ekleme 

mutation {
  insert_users(objects: [{
    email: "1@email.com",
    firstname: "name",
    lastname: "surname"
  }]) { 
    returning {
      id
      email
      firstname
      lastname
      created_at
      updated_at
    }
  }
}

3- Kullanıcı ve şirketleri birbirine bağlama 

mutation {
  insert_user_company(objects: {user_id: 1, company_id: 1}) {
    returning {
      id
      company {
        id
        name
        created_at
        updated_at
      }
      user {
        id
        email
        firstname
        lastname
        created_at
        updated_at
      }
    }
  }
}



4- Şirket silme

mutation {
  delete_companies_by_pk(id: 2) {
    id
  }
}


5- Kullanıcı silme 

mutation {
  delete_users_by_pk(id: 2) {
    id
  }
}


6- Şirket ve kullanıcı arasındaki bağı koparma 

mutation {
  delete_user_company(where: {_and: {user_id: {_eq: 1}, company_id: {_eq: 1}}}) {
    returning {
      company_id
      user_id
      id
    }
  }
}


7- Şirket düzenleme

mutation {
  update_companies_by_pk(
    pk_columns: {id: 1}, 
    _set: {
      name: "Company 1"
    }) {
    id
    name
    created_at
    updated_at
  }
}


8- Kullanıcı düzenleme

mutation {
  update_users_by_pk(
    pk_columns: {id: 1}, 
    _set: {
      email: "1@email.com", 
      firstname: "name updated", 
      lastname: "surname updated"
    }) {
    id
    email
    firstname
    lastname
    created_at
    updated_at
  }
}



9- Şirketi ve ona bağlı kullanıcıları listeleme

query MyQuery {
  companies_by_pk(id: 1) {
    id
    name
    created_at
    updated_at
    users {
      user {
        id
        email
        firstname
        lastname
        created_at
        updated_at
      }
    }
  }
}


10- Kullanıcıları ve onlara bağlı şirketleri listeleme.

query MyQuery {
  users_by_pk(id: 1) {
    id
    email
    firstname
    lastname
    created_at
    updated_at
    companies {
      company {
        id
        name
        created_at,
        updated_at
      }
    }
  }
}

query MyQuery {
  users {
    id
    email
    firstname
    lastname
    created_at
    updated_at
    companies {
      company {
        id
        name
        created_at
        updated_at
      }
    }
  }
}


```