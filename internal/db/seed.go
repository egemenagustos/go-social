package db

import (
	"context"
	"database/sql"
	"fmt"
	store "go-social/internal/storage"
	"log"
	"math/rand"

	"github.com/google/uuid"
)

var names = []string{
	"Ahmet",
	"Mehmet",
	"Mustafa",
	"Ali",
	"Hüseyin",
	"Hasan",
	"İbrahim",
	"İsmail",
	"Osman",
	"Yusuf",
	"Emre",
	"Burak",
	"Mert",
	"Can",
	"Kerem",
	"Onur",
	"Serkan",
	"Volkan",
	"Tolga",
	"Uğur",
	"Arda",
	"Eren",
	"Kaan",
	"Berk",
	"Batuhan",
	"Oğuz",
	"Sinan",
	"Furkan",
	"Barış",
	"Egemen",
	"Emir",
	"Enes",
	"Samet",
	"Harun",
	"Ömer",
	"Halil",
	"Murat",
	"Yasin",
	"Recep",
	"Cem",
	"Koray",
	"Levent",
	"Alper",
	"Alp",
	"Doğukan",
	"Gökhan",
	"Ferhat",
	"Tarık",
	"Salih",
	"Metin",
}

var titles = []string{
	"Go ile REST API Geliştirmeye Giriş",
	"Backend Geliştiriciler İçin PostgreSQL Rehberi",
	"Go'da Pointer Kullanımı Nasıl Çalışır?",
	"Microservice Mimarisi Nedir?",
	"Redis ile Uygulama Performansını Artırmak",
	"Docker ile Go Uygulaması Çalıştırmak",
	"REST API Tasarlarken Dikkat Edilmesi Gerekenler",
	"Go'da Goroutine ve Concurrency Mantığı",
	"PostgreSQL Index Kullanımı ve Performans",
	"JWT Authentication Nasıl Çalışır?",
	"Temiz Kod Yazmak İçin 10 İpucu",
	"Backend Projelerinde Katmanlı Mimari",
	"Go ile Database İşlemleri",
	"Redis Cache Kullanımı Ne Zaman Mantıklı?",
	"Docker Compose ile Development Ortamı Kurmak",
	"API Rate Limiting Nedir ve Nasıl Çalışır?",
	"Go Context Yapısını Anlamak",
	"SQL Sorgularında Performans Optimizasyonu",
	"Web Uygulamalarında Güvenlik Temelleri",
	"Monolith ve Microservice Arasındaki Farklar",
}

var contents = []string{
	"Go ile REST API geliştirirken routing, handler yapısı ve HTTP response yönetimi gibi temel konuları anlamak oldukça önemlidir.",
	"PostgreSQL güçlü ve güvenilir bir ilişkisel veritabanıdır. Doğru index kullanımı ve iyi tasarlanmış sorgular performansı ciddi şekilde artırabilir.",
	"Go'da pointer kullanımı bir değerin bellekteki adresine erişmemizi sağlar. Özellikle büyük struct yapılarında gereksiz kopyalamayı azaltmak için kullanılabilir.",
	"Microservice mimarisi büyük uygulamaları bağımsız servisler halinde geliştirmeyi sağlar. Ancak dağıtık sistem karmaşıklığını da beraberinde getirir.",
	"Redis genellikle cache, session yönetimi, distributed lock ve hızlı veri erişimi gereken senaryolarda kullanılan in-memory bir veri deposudur.",
	"Docker sayesinde Go uygulamalarını bağımlılıklarıyla birlikte paketleyebilir ve farklı ortamlarda aynı şekilde çalıştırabiliriz.",
	"İyi tasarlanmış bir REST API anlaşılır endpoint isimleri, doğru HTTP status kodları ve tutarlı response modelleri kullanmalıdır.",
	"Goroutine yapısı Go uygulamalarında aynı anda birden fazla işin verimli şekilde yürütülmesini sağlar ve concurrency işlemlerini oldukça kolaylaştırır.",
	"PostgreSQL indexleri doğru kolonlarda kullanıldığında sorgu sürelerini ciddi şekilde azaltabilir ancak gereksiz index kullanımı yazma performansını düşürebilir.",
	"JWT kullanıcı kimliğini doğrulamak ve API isteklerini yetkilendirmek için sıklıkla kullanılan token tabanlı bir authentication yöntemidir.",
	"Temiz kod yazarken anlamlı değişken isimleri kullanmak, küçük fonksiyonlar oluşturmak ve tekrar eden kodlardan kaçınmak önemlidir.",
	"Katmanlı mimari uygulamanın farklı sorumluluklarını birbirinden ayırarak kodun daha kolay test edilmesini ve bakım yapılmasını sağlar.",
	"Go'da database/sql paketi sayesinde SQL veritabanlarıyla doğrudan çalışabilir, sorgular çalıştırabilir ve transaction yönetebiliriz.",
	"Redis cache özellikle sık okunan ancak nadiren değişen verilerin hızlı bir şekilde kullanıcıya sunulması gereken durumlarda oldukça faydalıdır.",
	"Docker Compose birden fazla servisten oluşan uygulamaları tek bir yapılandırma dosyası üzerinden kolayca ayağa kaldırmamızı sağlar.",
	"Rate limiting bir API'nin belirli bir süre içerisinde kabul edeceği istek sayısını sınırlandırarak kötüye kullanım ve aşırı yüklenmeyi önler.",
	"Go context paketi request timeout, cancellation ve request bazlı değerlerin uygulamanın farklı katmanlarına taşınması için kullanılır.",
	"SQL performansını artırmak için gereksiz kolonları seçmemek, uygun indexler kullanmak ve execution planlarını incelemek önemlidir.",
	"Web uygulamalarında güvenlik için input validation, authentication, authorization ve güvenli veri saklama yöntemleri birlikte kullanılmalıdır.",
	"Monolith mimarisi başlangıçta daha basit bir geliştirme deneyimi sunarken microservice mimarisi bağımsız ölçeklenebilirlik ve deployment avantajı sağlar.",
}

var tags = []string{
	"go",
	"backend",
	"rest-api",
	"postgresql",
	"database",
	"redis",
	"docker",
	"devops",
	"microservices",
	"architecture",
	"concurrency",
	"goroutine",
	"pointer",
	"memory",
	"authentication",
	"jwt",
	"security",
	"clean-code",
	"performance",
	"http",
}

var comments = []string{
	"Gayet açıklayıcı bir yazı olmuş.",
	"Bu konuyu merak ediyordum, teşekkürler.",
	"Örnekler çok faydalı olmuş.",
	"Go öğrenirken işime yaradı.",
	"Devamını da bekliyorum.",
	"Oldukça sade ve anlaşılır anlatılmış.",
	"Bu konuyla ilgili daha fazla örnek güzel olur.",
	"Backend tarafında önemli bir konu.",
	"PostgreSQL kısmı özellikle faydalıydı.",
	"Performans tarafını güzel açıklamışsın.",
	"Ben de benzer bir yapı kullanıyorum.",
	"Gerçek proje örneğiyle anlatılması güzel olmuş.",
	"Bu detay daha önce kafamı karıştırıyordu.",
	"Teşekkürler, konu şimdi daha net.",
	"Microservice tarafında güzel bir özet olmuş.",
	"Redis kullanımını güzel açıklamışsın.",
	"Yeni başlayanlar için faydalı bir içerik.",
	"Bu yaklaşımı kendi projemde deneyeceğim.",
	"Gayet güzel bir kaynak olmuş.",
	"Bir sonraki yazıda test tarafını da anlatabilirsin.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete!")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		id, _ := uuid.NewV7()

		users[i] = &store.User{
			Id:       id.String(),
			Username: names[i%len(names)] + fmt.Sprintf("%d", i),
			Email:    names[i%len(names)] + fmt.Sprintf("%d", i) + "@mail.com",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		id, _ := uuid.NewV7()

		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			Id:      id.String(),
			Content: contents[rand.Intn(len(contents))],
			Title:   titles[rand.Intn(len(titles))],
			UserId:  user.Id,
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		id, _ := uuid.NewV7()
		cms[i] = &store.Comment{
			Id:      id.String(),
			PostId:  posts[rand.Intn(len(posts))].Id,
			UserId:  users[rand.Intn(len(users))].Id,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
